package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/nt2311-vn/snippetbox/internal/handlers"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	fileServer := http.FileServer(http.Dir(filepath.Join("ui", "static", "/")))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/snippet/view", handlers.SnippetView)
	mux.HandleFunc("/snippet/create", handlers.SnippetCreate)

	log.Printf("Starting server on %s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalln("error on starting server", err)
	}
}
