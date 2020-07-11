package validate

import "testing"

func TestIsValidEmail(t *testing.T) {
	type TestData struct {
		In  string
		Out bool
	}
	testDataList := []TestData{
		// valid emails
		{In: "alex@bell-labs.com", Out: true},
		{In: "rperl001@mit.edu", Out: true},
		// invalid emails
		{In: "", Out: false},
		{In: "yo!", Out: false},
	}
	for _, testData := range testDataList {
		if IsValidEmail(testData.In) != testData.Out {
			t.Errorf("expected %s to return %v but got %v", testData.In, testData.Out, !testData.Out)
		}
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	type TestData struct {
		In  string
		Out bool
	}
	testDataList := []TestData{
		// valid phone numbers
		{In: "+6139888998", Out: true},
		{In: "+61488224568", Out: true},
		// invalid phone numbers
		{In: "0488445688", Out: false},
		{In: "(03) 9333 7119", Out: false},
		{In: "+613AB88998", Out: false},
	}
	for _, testData := range testDataList {
		if IsValidPhoneNumber(testData.In) != testData.Out {
			t.Errorf("expected %s to return %v but got %v", testData.In, testData.Out, !testData.Out)
		}
	}
}
