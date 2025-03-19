package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

var (
	htmlStaticDir  = filepath.Join("ui", "html", "pages")
	htmlPartialDir = filepath.Join("ui", "html", "partials")
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		filepath.Join(htmlStaticDir, "base.html"),
		filepath.Join(htmlPartialDir, "nav.html"),
		filepath.Join(htmlStaticDir, "home.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err = ts.ExecuteTemplate(w, "base", nil); err != nil {
		log.Print(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
