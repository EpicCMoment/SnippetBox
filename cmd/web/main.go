package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	fileServerRoot := flag.String("fsroot", "./ui/static/", "Root folder of the file server")

	flag.Parse()

	fileServer := http.FileServer(http.Dir(*fileServerRoot))

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/getfile/", sendFile)

	log.Printf("Web server is started on localhost%s", *addr)

	err := http.ListenAndServe(*addr, mux)

	if err != nil {
		log.Fatal(err)
	}

}
