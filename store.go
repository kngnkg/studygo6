package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("postgres", "host=ch06-db user=gwp password=gwp dbname=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}

	// err = Db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
}

// get all posts
func Posts(limit int) (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, content, author FROM posts LIMIT $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

// Get a single post
func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("SELECT id, content, author FROM posts WHERE id = $1", id).
		Scan(&post.Id, &post.Content, &post.Author)
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	statement := "INSERT INTO posts (content, author) VALUES ($1, $2) RETURNING id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

// Update a post
func (post *Post) Update() (err error) {
	_, err = Db.Exec("UPDATE posts SET content = $2, author = $3 WHERE id = $1", post.Id, post.Content, post.Author)
	return
}

// Delete a post
func (post *Post) Delete() (err error) {
	_, err = Db.Exec("DELETE FROM posts WHERE id = $1", post.Id)
	return
}

// Delete all posts
func DeleteAll() (err error) {
	_, err = Db.Exec("DELETE FROM posts")
	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}

	// Create a post
	fmt.Println(post)
	post.Create()
	fmt.Println(post)

	// Get one post
	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost)

	// Update the post
	readPost.Content = "Bonjour Monde!"
	readPost.Author = "Pierre"
	readPost.Update()

	// Get all posts
	posts, _ := Posts(10)
	fmt.Println(posts)

	// Delete the post
	readPost.Delete()

	// Get all posts
	posts, _ = Posts(10)
	fmt.Println(posts)

}
