package main

import (
	"flag"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	fileServerRoot := flag.String("fsroot", "./ui/static/", "Root folder of the file server")
	fileServer := http.FileServer(http.Dir(*fileServerRoot))

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodGet, "/getfile/:filename", app.sendFile)

	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	middlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return middlewareChain.Then(router)
}