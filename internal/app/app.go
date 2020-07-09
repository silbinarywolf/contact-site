package app

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/silbinarywolf/contact-site/internal/config"
	"github.com/silbinarywolf/contact-site/internal/db"
)

const port = ":8080"

const databaseName = "ContactSite"

// templates are parsed once at boot-up so they only need to be parsed once and to
// catch any parsing problems as soon as possible.
//
// We store the files in ".assets" with a prefixed "." so that if we decide to serve
// our "static" files via Apache/Nginx, we can make the rules for public/privately exposed
// folders simple. (ie. all dot-prefixed folders are denied/blocked from public)
var templates = template.Must(template.ParseFiles(".assets/index.html"))

var (
	flagInit    bool
	flagDestroy bool
)

func init() {
	flag.BoolVar(&flagInit, "init", false, "if init flag is used, the database, tables and initial data will be setup")
	flag.BoolVar(&flagDestroy, "destroy", false, "if destroy flag is used, the database will be destroyed.")
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

type PhoneNumber struct {
	ID        int64
	ContactID int64
	Number    string
}

type Contact struct {
	ID           int64
	FullName     string
	Email        string
	PhoneNumbers []PhoneNumber
}

func mustDestroy() {
	db := db.Get()
	dropTables := []string{
		`DROP TABLE Contact`,
		`DROP TABLE PhoneNumber`,
	}
	for _, dropTableQuery := range dropTables {
		if _, err := db.Query(dropTableQuery); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P01" {
				// Do nothing if "undefined_table" error.
				// Just means table doesn't exist so if it never existed, thats fine.
			} else {
				panic(err)
			}
		}
	}
}

func mustSetupOrUpdate() {
	db := db.Get()

	// Create database
	/*_, err := db.Query("CREATE DATABASE " + databaseName + ";")
	if err != nil {
		return err
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			// Do nothing if "duplicate_database" error
			// Database has already been created.
		} else {
			panic(err)
		}
	}*/

	// Create tables
	createTables := []string{
		`CREATE TABLE PhoneNumber(
			ID        SERIAL PRIMARY KEY NOT NULL,
			ContactID INT              NOT NULL,
			Number    VARCHAR(16)      NOT NULL
		)`,
		`CREATE TABLE Contact(
			ID        SERIAL PRIMARY KEY NOT NULL,
			FullName  VARCHAR(255)     NOT NULL,
			Email     VARCHAR(255)     NOT NULL
		)`,
	}
	for _, createTableQuery := range createTables {
		if _, err := db.Query(createTableQuery); err != nil {
			panic(err)
			/*if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P07" {
				// Do nothing if "duplicate_table" error.
				// Just means table was already created.
			} else {
				panic(err)
			}*/
		}
	}
	// Fill with data
	records := []Contact{
		{
			FullName: "Alex Bell",
			Email:    "Fredrik Idestam",
			PhoneNumbers: []PhoneNumber{
				{Number: "03 8578 6688"},
				{Number: "1800728069"},
			},
		},
		{
			FullName: "Fredrik Idestam",
			PhoneNumbers: []PhoneNumber{
				{Number: "+6139888998"},
			},
		},
		{
			FullName: "Radia Perlman",
			Email:    "rperl001@mit.edu",
			PhoneNumbers: []PhoneNumber{
				{Number: "+6139888998"},
			},
		},
	}

	for _, record := range records {
		err := db.QueryRow(`INSERT INTO Contact (FullName, Email) VALUES ($1, $2) RETURNING ID`, record.FullName, record.Email).Scan(&record.ID)
		if err != nil {
			panic(err)
		}
		if record.ID == 0 {
			panic("Expected insertion to return an id not equal to 0")
		}
		for _, childRecord := range record.PhoneNumbers {
			err := db.QueryRow(`INSERT INTO PhoneNumber (ContactID, Number) VALUES($1, $2) RETURNING ID`, record.ID, childRecord.Number).Scan(&childRecord.ID)
			if err != nil {
				panic(err)
			}
			if childRecord.ID == 0 {
				panic("Expected insertion to return an id not equal to 0")
			}
		}
	}

	//db.Query("INSERT INTO Contact (ID, FullName, Email) VALUES (?, ?, ?)")
}

func Start() {
	flag.Parse()

	// Load config
	config.MustLoad()
	config := config.Get()

	// Connect the database
	db.Connect(db.Settings{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
	})
	defer db.Close()

	if flagDestroy {
		mustDestroy()
		os.Exit(0)
	}
	if flagInit {
		mustSetupOrUpdate()
		os.Exit(0)
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
