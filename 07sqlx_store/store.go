package main

// $ go install github.com/jmoiron/sqlx@latest
// $ go mod tidy

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Post struct {
	Id         int
	Content    string
	AuthorName string `db: author`
}

var Db *sqlx.DB

// connect to the Db
func init() {
	log.SetFlags(log.Llongfile)

	var err error
	Db, err = sqlx.Open("postgres", "host=ch06-db user=gwp password=gwp dbname=gwp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// デバッグ用
	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Get a single post
func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRowx("select id, content, author from posts where id = $1", id).
		StructScan(&post)
	if err != nil {
		log.Print(err)
		return
	}
	return
}

// Create a new post
func (post *Post) Create() (err error) {
	err = Db.QueryRow("insert into posts (content, author) values ($1, $2) returning id", post.Content, post.AuthorName).
		Scan(&post.Id)
	return
}

func main() {
	post := Post{Content: "Hello World!", AuthorName: "Sau Sheong"}
	post.Create()
	fmt.Println(post)
}
