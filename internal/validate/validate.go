package validate

import (
	"regexp"
)

// Validate email address
//
// Obviously not as complex as validating against the actual email address spec. but *probably* good enough for this.
//
// Sourced from: https://golangnews.org/2020/06/validating-an-email-address/
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validate E.164 phone number
//
// I considered dropping the optional + logic as per comment on StackOverflow below
// so that the data we have can "just work" with services like Twilio.
//
// I also considered just using the following library since it seems to allow for more precise validation
// and can convert most phone numbers to E.164 for you (which is nice for UX!) but went for the simplest thing
// to begin with. I can always improve this later. (https://github.com/dongri/phonenumber)
//
// Sourced from: https://stackoverflow.com/a/23299989/5013410
var phoneNumberRegex = regexp.MustCompile("^\\+?[1-9]\\d{1,14}$")

// IsValidEmail will check if the email is valid
//
// I considered making this also do a "net.LookupMX" check on the domain of the email
// but doing a potentially slow DNS lookup for validation seems a bit overkill.
func IsValidEmail(email string) bool {
	return len(email) >= 3 &&
		len(email) <= 255 &&
		emailRegex.MatchString(email)
}

// IsValidPhoneNumber will check if the phone number satisifes E.164 format
func IsValidPhoneNumber(phoneNumber string) bool {
	return len(phoneNumber) >= 1 &&
		len(phoneNumber) <= 15 &&
		phoneNumberRegex.MatchString(phoneNumber)
}
