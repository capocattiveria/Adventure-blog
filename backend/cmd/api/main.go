package main

import (
	"log"
	"net/http"

	"adventure-blog/internal/server"
)

func main() {
	srv := server.New()
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
