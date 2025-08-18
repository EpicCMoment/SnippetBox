package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.ariffil.com/internal/models"
)

type application struct {
	errorLog		*log.Logger
	infoLog			*log.Logger
	snippets		*models.SnippetModel
	templateCache	map[string]*template.Template
	formDecoder		*form.Decoder
}

func main() {

	errLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate | log.Ltime | log.Llongfile)
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate | log.Ltime | log.Lshortfile)

	infoLog.Println("Parsing the command line flags")
	// process command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	
	flag.Parse()

	// create database connection pool
	

	infoLog.Println("Creating a connection pool to database")

	db, err := openDB()

	if err != nil {
		errLog.Fatal(err.Error())
	}

	defer db.Close()


	infoLog.Println("Caching the HTML templates")

	templateCache, err := newTemplateCache()

	if err != nil {
		errLog.Fatal(err.Error())
	}

	formDecoder := form.NewDecoder()
	
	app := application {
		errorLog: errLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
	}


	srv := &http.Server{
		Addr: *addr,
		ErrorLog: app.errorLog,
		Handler: app.routes(),
	}
	
	app.infoLog.Printf("Web server is being started on localhost%s", *addr)

	err = srv.ListenAndServe()

	if err != nil {
		app.errorLog.Fatal(err)
	}

}

func openDB() (*sql.DB, error){

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPasswd := os.Getenv("MYSQL_PASSWORD")

	dataSourceString := fmt.Sprintf("%s:%s@/snippetbox?parseTime=true",
										mysqlUser,
										mysqlPasswd)	

	db, err := sql.Open("mysql", dataSourceString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil


}