package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog	*log.Logger
	infoLog		*log.Logger
}

func main() {

	app := application {
		errorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate | log.Ltime | log.Llongfile),
		infoLog: log.New(os.Stdout, "[INFO]\t", log.Ldate | log.Ltime | log.Lshortfile),
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	fileServerRoot := flag.String("fsroot", "./ui/static/", "Root folder of the file server")

	flag.Parse()

	fileServer := http.FileServer(http.Dir(*fileServerRoot))

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/getfile/", app.sendFile)

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: app.errorLog,
		Handler: mux,
	}
	
	app.infoLog.Printf("Web server is being started on localhost%s", *addr)

	err := srv.ListenAndServe()

	if err != nil {
		app.errorLog.Fatal(err)
	}

}
