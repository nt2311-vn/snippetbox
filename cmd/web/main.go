package main

import (
	"log"
	"net/http"

	"github.com/nt2311-vn/snippetbox/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/snippet/view", handlers.SnippetView)
	mux.HandleFunc("/snippet/create", handlers.SnippetCreate)

	log.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatalln("error on starting server", err)
	}
}
