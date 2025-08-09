package main

import (
	"flag"
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	fileServerRoot := flag.String("fsroot", "./ui/static/", "Root folder of the file server")
	fileServer := http.FileServer(http.Dir(*fileServerRoot))

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.HandleFunc("/getfile/", app.sendFile)

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}