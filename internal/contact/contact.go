package contact

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/silbinarywolf/contact-site/internal/db"
	"github.com/silbinarywolf/contact-site/internal/validate"
)

var (
	// Client-facing errors
	ErrInvalidFullName    = validate.NewError("Invalid Full Name provided. Name provided is too long.")
	ErrInvalidEmail       = validate.NewError("Invalid Email provided")
	ErrInvalidPhoneNumber = validate.NewError("Invalid Phone Number provided")

	// Internal (developer) errors
	errContactAlreadyExists     = errors.New("cannot insert Contact record that already exists")
	errPhoneNumberAlreadyExists = errors.New("cannot insert PhoneNumber record that already exists")
)

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

func InsertNew(record *Contact) error {
	db := db.Get()

	// Validate
	//
	// This could be in a new function, but that'd be premature as it'd only be called in one place, here.
	// So we leverage the amazing forgotten 70's technology of block-scoping
	// Aesthetically unappealing? Agreed. Practical and improves linear readability? Definitely.
	//
	// Also very easy to copy-paste to a new function when/if the time comes that I need this called
	// in two or more seperate places.
	{
		if record.ID != 0 {
			return errContactAlreadyExists
		}
		if len(record.FullName) >= 255 {
			return ErrInvalidFullName
		}
		// We allow a blank email address for these records
		// but that doesn't mean I want my email validation code to allow
		// blank strings, so we capture that information at this level
		if len(record.Email) != 0 &&
			!validate.IsValidEmail(record.Email) {
			return ErrInvalidEmail
		}
		for _, childRecord := range record.PhoneNumbers {
			if childRecord.ID != 0 {
				return errPhoneNumberAlreadyExists
			}
			// TODO(Jae): 2020-07-11
			// parse reasonable phone numbers into E.164 format
			if !validate.IsValidPhoneNumber(childRecord.Number) {
				return ErrInvalidPhoneNumber
			}
		}
	}

	// Insert record into DB
	err := db.QueryRow(`INSERT INTO Contact (FullName, Email) VALUES ($1, $2) RETURNING ID`, record.FullName, record.Email).Scan(&record.ID)
	if err != nil {
		return err
	}
	if record.ID == 0 {
		panic("Unexpected error. Failed get ID after inserting Contact record.")
	}
	for _, childRecord := range record.PhoneNumbers {
		err := db.QueryRow(`INSERT INTO PhoneNumber (ContactID, Number) VALUES($1, $2) RETURNING ID`, record.ID, childRecord.Number).Scan(&childRecord.ID)
		if err != nil {
			return err
		}
		if childRecord.ID == 0 {
			panic("Unexpected error. Failed get ID after inserting PhoneNumber record.")
		}
	}
	return nil
}

// MustInitialize will setup the necessary tables and add some mock data into the
// database.
//
// This function will panic if an error occurs.
func MustInitialize() {
	db := db.Get()

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
		}
	}
	// Fill with data
	records := []*Contact{
		{
			FullName: "Alex Bell",
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

	for i, record := range records {
		if err := InsertNew(record); err != nil {
			panic(fmt.Sprintf("Failed to insert record %d: %s", i, err))
		}
	}
}

func MustDestroy() {
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
