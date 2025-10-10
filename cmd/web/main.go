package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"snippetbox.ariffil.com/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	users         *models.UserModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	config        *StartupConfig
	sessManager   *scs.SessionManager
}

type StartupConfig struct {
	App struct {
		Port int `mapstructure: "port"`
		Host string `mapstructure: "host"`

		Database struct {
			Username string `mapstructure: "username"`
			Password string `mapstructure: "password"`
			Host     string `mapstructure: "host"`
			Port     int    `mapstructure: "port"`
		} `mapstructure: "database"`
	} `mapstructure: "app"`
}

func main() {

	// set the loggers
	errLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Llongfile)
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorLog: errLog,
		infoLog:  infoLog,
	}

	// Read in the startup config "config.yaml"
	stConf, err := getStartupConfig()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	app.config = stConf
	// Config done

	// Databse connection start
	infoLog.Println("Creating a connection pool to database")

	db, err := openDB(&app)

	if err != nil {
		errLog.Fatal(err.Error())
	}

	defer db.Close()

	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}
	// Database connection done

	// Create a session manager for use sessions
	app.sessManager = scs.New()
	app.sessManager.Store = mysqlstore.New(db)
	app.sessManager.Lifetime = 12 * time.Hour
	// Session creation done

	// Cache templating start

	infoLog.Println("Caching the HTML templates")

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err.Error())
	}

	app.templateCache = templateCache
	// Cache templating end

	// Set the form decoder
	formDecoder := form.NewDecoder()
	app.formDecoder = formDecoder
	// Form decoder setting done

	addr := fmt.Sprintf("%s:%d", stConf.App.Host, stConf.App.Port)

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Web server is being started on https://%s", addr)

	err = srv.ListenAndServeTLS("/home/rudrik/Desktop/SnippetBox/tls/cert.pem", "/home/rudrik/Desktop/SnippetBox/tls/key.pem")

	if err != nil {
		app.errorLog.Fatal(err)
	}

}

func getStartupConfig() (*StartupConfig, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/home/rudrik/Desktop/SnippetBox/")

	err := viper.ReadInConfig()
	if err != nil {
		cwd, _ := os.Getwd()
		return nil, fmt.Errorf("error reading the config file %s, %s", cwd+"/config.yaml", err)
	}

	sc := StartupConfig{}

	sc.App.Port = viper.GetInt("app.port")
	sc.App.Host = viper.GetString("app.host")
	sc.App.Database.Host = viper.GetString("app.database.host")
	sc.App.Database.Port = viper.GetInt("app.database.port")
	sc.App.Database.Username = viper.GetString("app.database.username")
	sc.App.Database.Password = viper.GetString("app.database.password")

	//app.infoLog.Printf("Startup config is read %+v", sc)

	return &sc, nil

}

func openDB(app *application) (*sql.DB, error) {

	mysqlUser := app.config.App.Database.Username
	mysqlPasswd := app.config.App.Database.Password

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
