package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.ariffil.com/ui"
)

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(http.FS(ui.Files))

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// middleware definitions

	// all handlers run after this middleware
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// this middleware provides session management for dynamically generated content
	dynamicMiddleware := alice.New(app.sessManager.LoadAndSave, noSurf, app.authenticate)

	// this middleware is used for endpoints requiring authentication
	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	// middleware definitions end

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	router.HandlerFunc(http.MethodGet, "/getfile/:filename", app.GET_sendFile)



	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.GET_home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddleware.ThenFunc(app.GET_snippetView))

	router.Handler(http.MethodGet, "/snippet/create", protectedMiddleware.ThenFunc(app.GET_snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protectedMiddleware.ThenFunc(app.POST_snippetCreate))

	router.Handler(http.MethodGet, "/user/login", dynamicMiddleware.ThenFunc(app.GET_login))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddleware.ThenFunc(app.POST_login))

	router.Handler(http.MethodGet, "/user/signup", dynamicMiddleware.ThenFunc(app.GET_signup))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddleware.ThenFunc(app.POST_signup))

	router.Handler(http.MethodPost, "/user/logout", protectedMiddleware.ThenFunc(app.POST_logout))

	return standardMiddleware.Then(router)
}
