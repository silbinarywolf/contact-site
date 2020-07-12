package app

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
	"github.com/silbinarywolf/contact-site/internal/config"
	"github.com/silbinarywolf/contact-site/internal/contact"
	"github.com/silbinarywolf/contact-site/internal/db"
	"github.com/silbinarywolf/contact-site/internal/validate"
)

var (
	flagInit    bool
	flagDestroy bool

	// templates holds all our /.templates files
	templates *template.Template

	isInitialized bool
	isClosed      bool
)

func init() {
	flag.BoolVar(&flagInit, "init", false, "if init flag is used, the database, tables and initial data will be setup")
	flag.BoolVar(&flagDestroy, "destroy", false, "if destroy flag is used, the database will be destroyed.")
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	type TemplateData struct {
		Contacts []contact.Contact
	}
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
	phoneNumbersDat := strings.TrimSpace(r.FormValue("PhoneNumbers"))
	var phoneNumbers []string
	if len(phoneNumbersDat) > 0 {
		if len(phoneNumbersDat) >= 4096 {
			// Arbitrarily limited the max amount of data to 4096.
			http.Error(w, "Invalid Phone Numbers given, too many phone numbers given.", http.StatusBadRequest)
			return
		}
		phoneNumbers = strings.Split(phoneNumbersDat, "\n")
	}

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

// MustInitialize will init various modules such as templates, configs and database connections.
//
// This logic was seperated from MustStart so that the initialization code could be blocking in our test code.
// The benefit of doing it this way is that it lowers the chance that the server may not have had enough time to start-up
// before our test code tries to make requests against our application.
func MustInitialize() {
	if isInitialized {
		panic("Cannot call MustInitialize more than once.")
	}

	// Initialize templates
	//
	// Templates are parsed once at boot-up so they only need to be parsed once and to
	// catch any parsing problems as soon as possible.
	//
	// We store the files in ".templates" with a prefixed "." so that if we decide to serve
	// our "static" files via Apache/Nginx, we can make the rules for public/privately exposed
	// folders simple. (ie. all dot-prefixed folders are denied/blocked from public)
	templates = template.Must(template.ParseFiles(
		".templates/index.html",
		".templates/postContact.html",
	))

	// Load config
	config.MustLoad()

	// Connect to the database
	config := config.Get()
	db.MustConnect(db.Settings{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
	})

	// Flags and initialization
	if flagDestroy {
		// --destroy flag will delete all tables
		mustDestroy()
		os.Exit(0)
	}
	if flagInit {
		// --init flag will only setup the tables/dummy data
		mustSetup()
		os.Exit(0)
	}
	mustSetup()

	// Setup routes
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

	isInitialized = true
}

func MustStart() {
	if !isInitialized {
		panic("Must call Initialize before calling Start")
	}
	port := ":" + strconv.Itoa(config.Get().Web.Port)
	log.Printf("Starting server on " + port + "...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

// MustClose should be called when the application closes.
func MustClose() {
	if isClosed {
		panic("Cannot call MustClose more than once.")
	}
	db.MustClose()
	isClosed = true
}

// mustDestroy will drop all the tables in the current database.
//
// In a real production situation, I'd probably make this hidden behind tag like "dev" or "debug"
// as it only exists for developer convenience.
func mustDestroy() {
	contact.MustDestroy()
}

// mustSetup will create tables and mock data for records.
func mustSetup() {
	contact.MustInitialize()
}
