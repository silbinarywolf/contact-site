package app

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
	"github.com/silbinarywolf/contact-site/internal/config"
	"github.com/silbinarywolf/contact-site/internal/contact"
	"github.com/silbinarywolf/contact-site/internal/db"
	"github.com/silbinarywolf/contact-site/internal/validate"
)

const port = ":8080"

// templates are parsed once at boot-up so they only need to be parsed once and to
// catch any parsing problems as soon as possible.
//
// We store the files in ".templates" with a prefixed "." so that if we decide to serve
// our "static" files via Apache/Nginx, we can make the rules for public/privately exposed
// folders simple. (ie. all dot-prefixed folders are denied/blocked from public)
var templates *template.Template

var (
	flagInit    bool
	flagDestroy bool
)

func init() {
	flag.BoolVar(&flagInit, "init", false, "if init flag is used, the database, tables and initial data will be setup")
	flag.BoolVar(&flagDestroy, "destroy", false, "if destroy flag is used, the database will be destroyed.")
}

type TemplateData struct {
	Contacts []contact.Contact
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	var templateData TemplateData
	templateData.Contacts = contact.GetAll()
	if err := templates.ExecuteTemplate(w, "index.html", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePostContact(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fullName := r.FormValue("FullName")
	email := r.FormValue("Email")
	phoneNumbersDat := r.FormValue("PhoneNumbers")
	phoneNumbersDat = strings.TrimSpace(phoneNumbersDat)
	if len(phoneNumbersDat) >= 4096 {
		// Arbitrarily limited the max amount of data to 4096.
		http.Error(w, "Invalid Phone Numbers given, too many phone numbers given.", http.StatusBadRequest)
		return
	}
	phoneNumbers := strings.Split(phoneNumbersDat, "\n")

	// Create record from request
	record := &contact.Contact{}
	record.FullName = fullName
	record.Email = email
	// We know the size of phone numbers provided.
	// So lets allocate precisely that amount.
	record.PhoneNumbers = make([]contact.PhoneNumber, len(phoneNumbers))
	for i, phoneNumber := range phoneNumbers {
		record.PhoneNumbers[i] = contact.PhoneNumber{
			Number: phoneNumber,
		}
	}
	if err := contact.InsertNew(record); err != nil {
		switch err := err.(type) {
		case *validate.ValidationError:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			log.Print(err)
			http.Error(w, "An unexpected error occurred inserting the record", http.StatusInternalServerError)
		}
		return
	}
	if err := templates.ExecuteTemplate(w, "postContact.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func Start() {
	// Initialize templates
	templates = template.Must(template.ParseFiles(
		".templates/index.html",
		".templates/postContact.html",
	))

	// Load config
	config.MustLoad()

	// Connect to the database
	config := config.Get()
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
		mustSetup()
		os.Exit(0)
	}

	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/postContact", handlePostContact)
	http.HandleFunc("/static/main.css", func(w http.ResponseWriter, r *http.Request) {
		// Manually serving CSS rather than using http.FileServer because Golang's in-built
		// detection methods can't really determine if the file is CSS or not.
		// Chrome complains if you try to load a CSS file with "text/plain". (has errors in Chrome DevTools)
		// See "DetectContentType" in the standard library, in file: net\http\sniff.go
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	log.Printf("Starting server on " + port + "...")
	//go func() {
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
	//}()
}

// mustDestroy will drop all the tables in the current database.
// In a real production situation, I'd probably make this hidden behind tag like "dev" or "debug"
// as it only exists for developer convenience.
func mustDestroy() {
	contact.MustDestroy()
}

// mustSetup will execute if the "flagInit" global variable is true
func mustSetup() {
	contact.MustInitialize()
}
