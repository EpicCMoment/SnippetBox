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

	// all handlers run after this middleware
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/getfile/:filename", app.sendFile)



	// this middleware provides session management for dynamically generated content
	dynamicMiddleware := alice.New(app.sessManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddleware.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamicMiddleware.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamicMiddleware.ThenFunc(app.snippetCreatePost))



	return standardMiddleware.Then(router)
}