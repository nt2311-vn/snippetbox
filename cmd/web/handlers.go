package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/nt2311-vn/snippetbox/internal/models"
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

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
