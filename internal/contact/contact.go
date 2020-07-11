package contact

import (
	"errors"

	"github.com/silbinarywolf/contact-site/internal/db"
	"github.com/silbinarywolf/contact-site/internal/validate"
)

var (
	ErrInvalidEmail       = errors.New("invalid Email provided")
	ErrInvalidPhoneNumber = errors.New("invalid Phone Number provided")

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

func InsertNew(record Contact) error {
	db := db.Get()

	// Validate
	//
	// This could be in a new function, but that'd be premature as it'd only
	// be called in one place, here. So we leverage the amazing forgotten 70's technology of block-scoping
	// Aesthetically unappealing? Sure. Practical and improves linear readability? Yes.
	{
		if record.ID == 0 {
			return errContactAlreadyExists
		}
		if !validate.IsValidEmail(record.Email) {
			return ErrInvalidEmail
		}
		for _, childRecord := range record.PhoneNumbers {
			if childRecord.ID == 0 {
				return errPhoneNumberAlreadyExists
			}
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
