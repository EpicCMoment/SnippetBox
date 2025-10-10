package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)

}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	templateCheckBuffer := bytes.Buffer{}

	err := ts.ExecuteTemplate(&templateCheckBuffer, "base", data)

	if err != nil {
		app.serverError(w, err)
		return
	}

	writeCount, _ := templateCheckBuffer.WriteTo(w)

	if writeCount < int64(templateCheckBuffer.Len()) {
		app.serverError(w, errors.New("dynamic template generation error"))
		return
	}

}

func (a *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       a.sessManager.PopString(r.Context(), "flash"),
	}
}
