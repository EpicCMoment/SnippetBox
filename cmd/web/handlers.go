package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

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

	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expiresString := r.PostForm.Get("expires")
	expires, err := strconv.Atoi(expiresString)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank!"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long!"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must be 1, 7 or 365!"
	} 

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}


	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	newSnippetURL := fmt.Sprintf("/snippet/view/%d", id)

	http.Redirect(w, r, newSnippetURL, http.StatusSeeOther)

}
