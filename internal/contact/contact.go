package contact

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/nyaruka/phonenumbers"

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

func InsertNew(record *Contact) (rErr error) {
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
		for i, _ := range record.PhoneNumbers {
			childRecord := &record.PhoneNumbers[i]
			if childRecord.ID != 0 {
				return errPhoneNumberAlreadyExists
			}
			phoneNumber := strings.TrimSpace(childRecord.Number)
			// Validate phone number against Australian format as the test data provided to me
			// implied that we should infer Australian numbers.
			//
			// I initially stumbled across this parsing/formatting implementation: https://github.com/dongri/phonenumber
			// but it didn't fill me with much confidence asE.164 is seemingly like timezones, wherein they change
			// requirements over time. I ideally want to buy-in to something that is maintained or easy to take over maintenance for.
			//
			// So then I discovered that Google had libraries dedicated to parsing this but only C/Java/JavaScript implementations:
			// - https://github.com/google/libphonenumber
			//
			// So finally, after more googling I lucked upon this Golang implementation based on Google's Java implementation. 
			// It has reasonable tests and instructions on how to update the binary data. Promising! So I'm rolling with it.
			// - https://github.com/nyaruka/phonenumbers 
			parsedNumber, err := phonenumbers.Parse(phoneNumber, "AU")
			if err != nil {
				return ErrInvalidPhoneNumber
			}
			formattedNum := phonenumbers.Format(parsedNumber, phonenumbers.E164)

			// It feels like a bit of a code smell for the validation of this record
			// to modify the phone numbers. But seems to be the best spot
			// to put this logic for now, so, I'll just do it. If I get a better idea
			// on where to place this, I'll can always move it later.
			childRecord.Number = formattedNum
		}
	}

	// Insert record into DB
	tx, err := db.Get().Begin()
	if err != nil {
		return err
	}
	hasCommitted := false
	defer func(){
		if hasCommitted {
			return
		}
		if err := tx.Rollback(); err != nil {
			rErr = err
		}
	}()
	err = tx.QueryRow(`INSERT INTO Contact (FullName, Email) VALUES ($1, $2) RETURNING ID`, record.FullName, record.Email).Scan(&record.ID)
	if err != nil {
		return err
	}
	if record.ID == 0 {
		panic("Unexpected error. Failed get ID after inserting Contact record.")
	}
	for _, childRecord := range record.PhoneNumbers {
		err := tx.QueryRow(`INSERT INTO PhoneNumber (ContactID, Number) VALUES($1, $2) RETURNING ID`, record.ID, childRecord.Number).Scan(&childRecord.ID)
		if err != nil {
			return err
		}
		if childRecord.ID == 0 {
			panic("Unexpected error. Failed get ID after inserting PhoneNumber record.")
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	hasCommitted = true
	return
}

func GetAll() []Contact {
	db := db.Get()

	// I considered using an INNER JOIN like this:
	// - INNER JOIN PhoneNumber ON PhoneNumber.ContactID = Contact.ID
	// But ultimately just opted to do a query per records has_many for simplicity
	// and easier extensibility. (ie. adding more relationships, etc)
	rows, err := db.Query(`SELECT ID, FullName, Email FROM Contact`)
	if err != nil {
		panic(err)
	}
	var contacts []Contact
	for rows.Next() {
		record := Contact{}
		err := rows.Scan(&record.ID, &record.FullName, &record.Email)
		if err != nil {
			panic(err)
		}
		childRows, err := db.Query(`SELECT ID, ContactID, Number FROM PhoneNumber WHERE ContactID = $1`, record.ID)
		if err != nil {
			panic(err)
		}
		for childRows.Next() {
			childRecord := PhoneNumber{}
			err := childRows.Scan(&childRecord.ID, &childRecord.ContactID, &childRecord.Number)
			if err != nil {
				panic(err)
			}
			record.PhoneNumbers = append(record.PhoneNumbers, childRecord)
		}
		contacts = append(contacts, record)
	}
	return contacts
}

// MustInitialize will setup the necessary tables and add some mock data into the
// database.
//
// This function will panic if an error occurs.
func MustInitialize() {
	db := db.Get()

	// Create tables
	createTables := []string{
		`CREATE TABLE Contact(
			ID        SERIAL PRIMARY KEY NOT NULL,
			FullName  VARCHAR(255)     NOT NULL,
			Email     VARCHAR(255)     NOT NULL
		)`,
		`CREATE TABLE PhoneNumber(
			ID        SERIAL PRIMARY KEY NOT NULL,
			ContactID INT              NOT NULL,
			Number    VARCHAR(16)      NOT NULL,
			CONSTRAINT FkContactID FOREIGN KEY (ContactID) REFERENCES Contact (ID)
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
				{Number: "(03) 9333 7119"},
				{Number: "0488445688"},
				{Number: "+61488224568"},
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
		`DROP TABLE PhoneNumber`,
		`DROP TABLE Contact`,
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
