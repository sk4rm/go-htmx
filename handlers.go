package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// GET /posts/<id>
	//  Shows post with matching id
	case http.MethodGet:
		id := r.URL.Path[len("/posts/"):]
		log.Println("GET /posts/" + id)

		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		check(err)
		defer conn.Close(context.Background())

		query := "SELECT * FROM posts WHERE id=" + id
		var post_id, title, body string
		err = conn.QueryRow(context.Background(), query).Scan(&post_id, &title, &body)

		if err == pgx.ErrNoRows {
			fmt.Fprintln(w, "no post with id "+id)
		} else {
			check(err)
			fmt.Fprintf(w, "post %s:\ntitle: %s\nbody: %s", post_id, title, body)
		}

	// POST /posts/
	//  Creates a new post and returns a unique id
	case http.MethodPost:
		log.Println("POST /posts/")

		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		check(err)
		defer conn.Close(context.Background())

		query := "INSERT INTO public.posts (title, body) VALUES ($1::text, $2::text) returning id;"
		var post_id string
		err = conn.QueryRow(context.Background(), query, title, body).Scan(&post_id)
		check(err)

		log.Printf("post #%s successfully created\n", post_id)
		w.Header().Set("HX-Redirect", "/posts/"+post_id)
		fmt.Fprintf(w, "post #%s successfully created\n", post_id)
	}
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/new-post.html")
	check(err)
	tmpl.Execute(w, nil)
	check(err)
}
