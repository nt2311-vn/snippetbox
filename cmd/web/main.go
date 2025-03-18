package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatalln("error on starting server", err)
	}
}

func home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello from snippet box"))
}

func snippetView(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}
