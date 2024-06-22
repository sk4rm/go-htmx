package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DATABASE_URL string = "db.sqlite"
	PORT         int    = 8080
)

var db *sql.DB

func main() {
	if DATABASE_URL == "" {
		log.Fatal("database url not provided")
	}

	var err error
	db, err = sql.Open("sqlite3", "db.sqlite")
	check(err)
	log.Printf("Established connection to %v", DATABASE_URL)
	defer db.Close()

	// Routes
	http.HandleFunc("GET /posts/{id}", viewPostHandler)
	http.HandleFunc("GET /posts/new/", newPostHandler)
	http.HandleFunc("GET /posts/", viewAllPostsHandler)
	http.HandleFunc("POST /posts/", postHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
