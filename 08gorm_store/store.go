package main

// $ go install github.com/jinzhu/gorm@latest
// $ go mod tidy

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Post struct {
	Id        int
	Content   string
	Author    string `sql:"not null"`
	Comments  []Comment
	CreatedAt time.Time
}

type Comment struct {
	Id        int
	Content   string
	Author    string `sql:"not null"`
	PostId    int
	CreatedAt time.Time
}

var Db *gorm.DB

// connect to the DB
func init() {
	log.SetFlags(log.Llongfile)

	var err error
	Db, err = gorm.Open("postgres", "host=ch06-db user=gwp password=gwp dbname=gwp sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	Db.AutoMigrate(&Post{}, &Comment{})
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}
	fmt.Println(post)

	// Create a post
	Db.Create(&post)
	fmt.Println(post)

	// Add a comment
	comment := Comment{Content: "Good post!", Author: "Joe"}
	Db.Model(&post).Association("Comments").Append(comment)

	// Get comments from a post
	var readPost Post
	Db.Where("author = $1", "Sau Sheong").First(&readPost)
	var comments []Comment
	Db.Model(&readPost).Related(&comments)
	fmt.Println(comments[0])
}
