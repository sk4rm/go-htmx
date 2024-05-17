package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// POST /posts/
//
// Creates a new post and returns a newly generated id.
func postHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /posts/")

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	check(err)
	defer conn.Close(context.Background())

	query := "INSERT INTO public.posts (title, body) VALUES ($1::text, $2::text) returning id;"
	var postID string

	// Create new record in database.
	err = conn.QueryRow(context.Background(), query, title, body).Scan(&postID)
	check(err)

	log.Printf("created post #%s\n", postID)

	w.Header().Set("HX-Redirect", "/posts/"+postID)
}

// GET /posts/{id}
//
// Shows post with matching id.
func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	log.Println("GET /posts/" + id)

	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "post id must only contain numbers", http.StatusBadRequest)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	check(err)

	query := "SELECT * FROM posts WHERE id=$1"
	var postID, title, body string
	err = conn.QueryRow(context.Background(), query, id).Scan(&postID, &title, &body)

	if err == pgx.ErrNoRows {
		fmt.Fprintln(w, "no post with id "+id)
	} else {
		check(err)
		fmt.Fprintf(w, "post %s:\ntitle: %s\nbody: %s", postID, title, body)
	}
}

// GET /posts/new/
//
// Displays frontend for creating new posts.
func newPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /posts/new/")

	tmpl, err := template.ParseFiles("templates/base.html", "templates/new-post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /posts/
//
// Display all posts
func viewAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /posts/")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	check(err)
	defer conn.Close(context.Background())

	query := "SELECT * FROM posts"
	rows, err := conn.Query(context.Background(), query)
	check(err)
	defer rows.Close()

	isEmpty := true
	for rows.Next() {
		isEmpty = false

		var postID, title, body string
		err := rows.Scan(&postID, &title, &body)
		check(err)

		_, err = fmt.Fprintf(w, "<a href=\"/posts/%v\">%v</a><br><br>", postID, title)
		check(err)
	}
	check(rows.Err())

	if isEmpty {
		fmt.Fprintln(w, "no posts :<")
		fmt.Fprintln(w, "make one at /posts/new? :D")
	}
}
