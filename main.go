package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatalln("database url not provided")
	}

	// Test database connection.
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	check(err)
	defer conn.Close(context.Background())

	// https://go.dev/blog/routing-enhancements
	http.HandleFunc("GET /posts/{id}", viewPostHandler)
	http.HandleFunc("GET /posts/new/", newPostHandler)
	http.HandleFunc("GET /posts/", viewAllPostsHandler)
	http.HandleFunc("POST /posts/", postHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
