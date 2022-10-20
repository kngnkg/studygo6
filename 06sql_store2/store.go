package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Post struct {
	Id       int
	Content  string
	Author   string
	Comments []Comment
}

type Comment struct {
	Id      int
	Content string
	Author  string
	Post    *Post
}

var Db *sql.DB

func init() {
	// ログの設定
	log.SetPrefix("[DEBUG]")
	log.SetFlags(log.Llongfile)

	var err error
	Db, err = sql.Open("postgres", "host=ch06-db user=gwp password=gwp dbname=gwp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// デバッグ用
	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("投稿が見つかりません")
		return
	}
	err = Db.QueryRow("INSERT INTO comments (content, author, post_id) VALUES ($1, $2, $3) RETURNING id", comment.Content, comment.Author, comment.Post.Id).
		Scan(&comment.Id)
	return
}

// Get a single post
func GetPost(id int) (post Post, err error) {
	post = Post{}
	post.Comments = []Comment{}
	err = Db.QueryRow("SELECT id, content, author FROM posts WHERE id = $1", id).
		Scan(&post.Id, &post.Content, &post.Author)
	if err != nil {
		log.Print(err)
		return
	}

	rows, err := Db.Query("SELECT id, content, author FROM comments WHERE post_id = $1", id)
	if err != nil {
		log.Print(err)
		return
	}
	for rows.Next() {
		comment := Comment{Post: &post}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			log.Print(err)
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	rows.Close()
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	err = Db.QueryRow("INSERT INTO posts (content, author) VALUES ($1, $2) RETURNING id", post.Content, post.Author).
		Scan(&post.Id)
	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	post.Create()
	fmt.Println(post)

	// Add a comment
	comment := Comment{Content: "Good post!", Author: "Joe", Post: &post}
	comment.Create()
	readPost, _ := GetPost(post.Id)

	fmt.Println(readPost)
	fmt.Println(readPost.Comments)
	fmt.Println(readPost.Comments[0].Post)
}
