package main

import (
	"log"
	"net/http"
	"os"

	"adventure-blog/internal/db"
	"adventure-blog/internal/server"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://adventure:adventure_secret@postgres:5432/adventure_blog?sslmode=disable"
	}

	pool, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	srv := server.New(pool)
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
