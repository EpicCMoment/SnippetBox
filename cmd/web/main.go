package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)


var infoLog *log.Logger = log.New(os.Stdout, "[INFO]\t", log.Ldate | log.Ltime | log.Lshortfile)
var errorLog *log.Logger = log.New(os.Stderr, "[ERROR]\t", log.Ldate | log.Ltime | log.Llongfile)

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

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}
	
	infoLog.Printf("Web server is being started on localhost%s", *addr)

	err := srv.ListenAndServe()

	if err != nil {
		errorLog.Fatal(err)
	}

}
