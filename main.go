package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	// Test database connection.
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	check(err)
	defer conn.Close(context.Background())

	http.HandleFunc("/posts/", postHandler)
	http.HandleFunc("/posts/new/", newPostHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
