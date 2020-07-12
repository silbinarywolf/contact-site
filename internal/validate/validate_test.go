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
