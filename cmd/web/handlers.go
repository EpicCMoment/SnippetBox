package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func (app *application) sendFile(w http.ResponseWriter, r *http.Request) {

	app.infoLog.Println("entered the sendFile handler")

	app.infoLog.Printf("getfile: %s\n", r.URL.Path)

	requestedFile := strings.TrimPrefix(r.URL.Path, "/getfile/")

	app.infoLog.Printf("requested file is: %s\n", requestedFile)

	http.ServeFile(w, r, filepath.Clean((requestedFile)))

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	templateFiles := []string{
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	ts, err := template.ParseFiles(templateFiles...)

	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	snippetIdString := r.URL.Query().Get("id")

	snippetId, err := strconv.Atoi(snippetIdString)

	if err != nil || snippetId < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Displaying the snippet %d", snippetId)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Create a new snippet...")

}
