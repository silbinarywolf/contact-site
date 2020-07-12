package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/silbinarywolf/contact-site/internal/app"
	"github.com/silbinarywolf/contact-site/internal/config"
)

var (
	HostName string
)

// TestMain will execute before all tests and allows us to do setup/teardown
func TestMain(m *testing.M) {
	// If we cannot find the config file in the current directory,
	// change cwd to 1 level above. This allows "go test" to be run
	// naively without building a specific Go test binary and placing
	// it in the correct directory.
	if !config.Exists() {
		dir, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("failed to get current dir: %s", err))
		}
		if err := os.Chdir(filepath.Join(dir, "..")); err != nil {
			panic(fmt.Sprintf("failed to change dir: %s", err))
		}
		// Fallthrough, if the config file still doesn't exist
		// MustInitialize will fail with the appropriate error message.
	}

	// Initialize the app
	app.MustInitialize()
	defer app.MustClose()

	// Set hostname we hit with get/post requests in our tests below
	HostName = "http://localhost:" + strconv.Itoa(config.Get().Web.Port)

	// Start application without blocking (so we can run tests)
	go app.MustStart()

	// Runs all the Test*** functions
	os.Exit(m.Run())
}

func TestGetHomePage(t *testing.T) {
	resp, err := http.Get(HostName)
	if err != nil {
		t.Fatalf("post error: path \"/\": %s", err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("readAll error: %s", err)
	}
	// This was here for debug purposes when using the -v verbose flag
	// t.Logf("%s", dat)
}

func TestPostFormSuccess(t *testing.T) {
	// Opted to just post data directly to the web server. Seemed like the most
	// robust way to test whether the server is running correctly or not.
	// Slow? Probably. But if it turns out to not be a good idea, we can always change it
	// later.
	resp, err := http.PostForm(
		HostName+"/postContact",
		url.Values{
			"FullName":     {"Test"},
			"Email":        {"test@test.com"},
			"PhoneNumbers": {"043"},
		},
	)
	if err != nil {
		t.Fatalf(
			"post error: path \"%s/\": %s",
			"postContact",
			err,
		)
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("readAll error: %s", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		// expected result, success!
	case http.StatusBadRequest:
		t.Errorf("unexpected response: %s", dat)
	default:
		t.Fatalf("unhandled error: %s", err)
	}
}

func TestPostFormFailure(t *testing.T) {
	resp, err := http.PostForm(
		HostName+"/postContact",
		url.Values{
			"FullName":     {"Test"},
			"Email":        {"BAD_EMAIL_TO_FAIL_VALIDATION"},
			"PhoneNumbers": {"043"},
		},
	)
	if err != nil {
		t.Fatalf(
			"post error: path \"%s/\": %s",
			"postContact",
			err,
		)
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("readAll error: %s", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		t.Errorf("unexpected response: %s", dat)
	case http.StatusBadRequest:
		// expected result, success!
		// If we add more failure cases, we might want to have an error code
		// attached to the header or something so we can detect the type of
		// error without needing to do string matches.
	default:
		t.Fatalf("unhandled error: %s", err)
	}
}
