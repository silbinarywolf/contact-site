package app

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

const port = ":8080"

// templates are parsed once at boot-up so they only need to be parsed once and to
// catch any parsing problems as soon as possible.
//
// We store the files in ".assets" with a prefixed "." so that if we decide to serve
// our "static" files via Apache/Nginx, we can make the rules for public/privately exposed
// folders simple. (ie. all dot-prefixed folders are denied/blocked from public)
var templates = template.Must(template.ParseFiles(".assets/index.html"))

type TemplateData struct {
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	var p TemplateData
	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Contact struct {
}

func Start() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	//dbHost = "192.168.99.100"
	//dbHost = "db"
	//dbUser = "admin"
	//dbPass = "password"
	//dbPort = "5432"

	shouldEarlyExit := false
	if dbUser == "" {
		log.Printf("DB_USER environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if dbPass == "" {
		log.Printf("DB_PASSWORD environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if dbPort == "" {
		log.Printf("DB_PORT environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if shouldEarlyExit {
		os.Exit(1)
	}
	// log.Printf("DEBUG: User: %v, Pass: %v, Port: %v\n", dbUser, dbPass, dbPort)
	db := pg.Connect(&pg.Options{
		Addr:     dbHost + ":" + dbPort,
		User:     dbUser,
		Password: dbPass,
		OnConnect: func(cn *pg.Conn) error {
			log.Printf("Successfully connected to SQL server")
			return nil
		},
	})
	for i := 0; i < 5; i++ {
		var n int
		_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
		if err == nil {
			break
		}
		log.Printf("Database connection attempt #%d: %v\n", i, err)
		if i == 4 {
			log.Println("Unable to connect to database. Stopping app.")
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}
	log.Println("Database connection successful")
	models := []interface{}{
		(*Contact)(nil),
	}
	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true, // temp table
		})
		if err != nil {
			panic(err)
		}
	}
	defer db.Close()

	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/static/main.css", func(w http.ResponseWriter, r *http.Request) {
		// Manually serving CSS rather than using http.FileServer because Golang's in-built
		// detection methods can't really determine if the file is CSS or not.
		// Chrome complains if you try to load a CSS file with "text/plain". (has errors in Chrome DevTools)
		// See "DetectContentType" in the standard library, in file: net\http\sniff.go
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	log.Printf("Starting server on " + port + "...\n")
	http.ListenAndServe(port, nil)
}
