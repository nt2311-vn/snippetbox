package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

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

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n\n", snippet)
	}

	// files := []string{
	// 	filepath.Join(htmlStaticDir, "base.html"),
	// 	filepath.Join(htmlPartialDir, "nav.html"),
	// 	filepath.Join(htmlStaticDir, "home.html"),
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", "text/html")
	//
	// if err = ts.ExecuteTemplate(w, "base", nil); err != nil {
	// 	log.Print(err.Error())
	// 	app.serverError(w, err)
	// }
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

	files := []string{
		filepath.Join(htmlStaticDir, "base.html"),
		filepath.Join(htmlPartialDir, "nav.html"),
		filepath.Join(htmlStaticDir, "view.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err = ts.ExecuteTemplate(w, "base", snippet); err != nil {
		app.serverError(w, err)
	}
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
