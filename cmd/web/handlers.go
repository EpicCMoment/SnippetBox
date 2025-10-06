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
	"snippetbox.ariffil.com/internal/validator"
)

type snippetCreateForm struct {
	Title		string	`form:"title"`
	Content		string	`form:"content"`
	Expires		int		`form:"expires"`
	validator.Validator	`form:"-"`
}

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

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}


	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.BelowMaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(int(form.Expires), 1, 7, 365), "expires", "This field should be 1, 7 or 365")

	if !form.Valid() {

		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)

		return
	}


	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessManager.Put(r.Context(), "flash", "Snippet successfully created!")



	newSnippetURL := fmt.Sprintf("/snippet/view/%d", id)
	http.Redirect(w, r, newSnippetURL, http.StatusSeeOther)

}
