package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const port = ":8080"

// templates are parsed once at boot-up so they only need to be parsed once and to
// catch any parsing problems as soon as possible.
//
// We store the files in ".assets" with a prefixed "." so that if we decide to serve
// our "static" files via Apache/Nginx, we can make the rules for public/privately exposed
// folders simple. (ie. all dot-prefixed folders are denied/blocked from public)
var templates = template.Must(template.ParseFiles(".assets/index.html"))

var (
	flagInit bool
)

func init() {
	flag.BoolVar(&flagInit, "init", false, "if init flag is used, the database, tables and initial data will be setup")
}

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

type Config struct {
	Database struct {
		Server   string `json:"server,omitempty"`
		Port     int    `json:"port,omitempty"`
		User     string `json:"user,omitempty"`
		Password string `json:"password,omitempty"`
	}
}

type PhoneNumber struct {
	ID     int
	Number string
}

type Contact struct {
	ID           int
	FullName     string
	Email        string
	PhoneNumbers []PhoneNumber
}

func loadConfig() Config {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	dat, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Config load error: %s\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(dat, &config); err != nil {
		log.Printf("Config parse error: %s\n", err)
		os.Exit(1)
	}
	return config
}

func setup() {
	_, err = db.Query("CREATE DATABASE " + dbName + ";")
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			// Do nothing if "duplicate_database" error
			// it's already been created
		} else {
			panic(err)
		}
	}
	_, err = db.Query(`
		DROP TABLE PhoneNumber;
		DROP TABLE Contact;
	`)
	_, err = db.Query(`
		CREATE TABLE PhoneNumber(
			ID        INT PRIMARY KEY  NOT NULL,
			ContactID INT              NOT NULL,
			Number    VARCHAR(16)      NOT NULL
		);
		CREATE TABLE Contact(
			ID        INT PRIMARY KEY  NOT NULL,
			FullName  VARCHAR(255)     NOT NULL,
			Email     VARCHAR(255)     NOT NULL
		);
	`)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P07" {
			// Do nothing if "duplicate_table" error, we're already initialized
		} else {
			panic(err)
		}
	}
}

func Start() {
	// Load config file
	config := loadConfig()

	dbUser := config.Database.User
	dbPass := config.Database.Password
	dbHost := config.Database.Server
	dbPort := strconv.Itoa(config.Database.Port)
	dbName := "ContactName"

	shouldEarlyExit := false
	if dbUser == "" {
		log.Printf("DB_USER environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if dbPass == "" {
		log.Printf("DB_PASSWORD environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if dbHost == "" {
		log.Printf("DB_HOST environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if dbPort == "" {
		log.Printf("DB_PORT environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if shouldEarlyExit {
		os.Exit(1)
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPass,
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Try and connect to the database
	for i := 0; i < 5; i++ {
		err := db.Ping()
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

	if flagInit {
		setup()
		os.Exit(0)
	}

	// Init
	_, err = db.Query("SELECT DATABASE " + dbName + ";")
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			// Do nothing if "duplicate_database" error
			// it's already been created
		} else {
			panic(err)
		}
	}

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
