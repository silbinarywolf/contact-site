package validate

import (
	"regexp"
)

// ValidationError is a distinct error type that we use when we want to expose
// error information to the frontend / end-user.
type ValidationError struct {
	message string
}

// assert at compile-time that this type satisfies the error interface
var _ error = new(ValidationError)

func (err *ValidationError) Error() string {
	return err.message
}

func NewError(message string) *ValidationError {
	return &ValidationError{
		message: message,
	}
}

// Validate email address
//
// Obviously not as complex as validating against the actual email address spec. but *probably* good enough for this.
//
// Sourced from: https://golangnews.org/2020/06/validating-an-email-address/
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsValidEmail will check if the email is valid
//
// I considered making this also do a "net.LookupMX" check on the domain of the email
// but doing a potentially slow DNS lookup for validation seems a bit overkill.
func IsValidEmail(email string) bool {
	return len(email) >= 3 &&
		len(email) <= 255 &&
		emailRegex.MatchString(email)
}
