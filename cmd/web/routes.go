package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() http.Handler {
	fileServer := http.FileServer(neuteredFileSystem{
		http.Dir(filepath.Join("ui", "static", "/")),
	})

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return secureHeaders(mux)
}
