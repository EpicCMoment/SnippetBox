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
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// middleware definitions

	// all handlers run after this middleware
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// this middleware provides session management for dynamically generated content
	dynamicMiddleware := alice.New(app.sessManager.LoadAndSave, noSurf)

	// this middleware is used for endpoints requiring authentication
	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	// middleware definitions end

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/getfile/:filename", app.sendFile)



	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddleware.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/snippet/create", protectedMiddleware.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protectedMiddleware.ThenFunc(app.snippetCreatePost))

	router.Handler(http.MethodGet, "/user/login", dynamicMiddleware.ThenFunc(app.serveLoginPage))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddleware.ThenFunc(app.login))

	router.Handler(http.MethodGet, "/user/signup", dynamicMiddleware.ThenFunc(app.serveSignupPage))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddleware.ThenFunc(app.signup))

	router.Handler(http.MethodPost, "/user/logout", protectedMiddleware.ThenFunc(app.logout))

	return standardMiddleware.Then(router)
}
