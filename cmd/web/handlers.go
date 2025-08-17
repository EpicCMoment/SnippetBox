package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"snippetbox.ariffil.com/internal/models"
)

func (app *application) sendFile(w http.ResponseWriter, r *http.Request) {

	app.infoLog.Println("entered the sendFile handler")

	app.infoLog.Printf("getfile: %s\n", r.URL.Path)

	requestedFile := strings.TrimPrefix(r.URL.Path, "/getfile/")

	app.infoLog.Printf("requested file is: %s\n", requestedFile)

	http.ServeFile(w, r, filepath.Clean((requestedFile)))

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	latestSnippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	data.Snippets = latestSnippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)


}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	snippetId, err := strconv.Atoi(params.ByName("id"))

	if err != nil || snippetId < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(snippetId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)


}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Display the form for creating a new snippet...")

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	newSnippetURL := fmt.Sprintf("/snippet/view/%d", id)

	http.Redirect(w, r, newSnippetURL, http.StatusSeeOther)

}
