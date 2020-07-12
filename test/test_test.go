package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/silbinarywolf/contact-site/internal/app"
	"github.com/silbinarywolf/contact-site/internal/config"
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

	// Start application without blocking (so we can run tests)
	go app.MustStart()

	// Runs all the Test*** functions
	os.Exit(m.Run())
}

func TestGetHomePage(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		t.Fatalf("post error: path \"/\": %s", err)
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("readAll error: %s", err)
	}
	t.Logf("%s", dat)
}

func TestPostFormSuccess(t *testing.T) {
	type TestData struct {
		In  url.Values
		Out bool
	}

	// Opted to just post data directly to the web server. Seemed like the most
	// effective way to test whether the server is running correctly or not.
	resp, err := http.PostForm(
		"http://127.0.0.1:8080/postContact",
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
		// success
	case http.StatusBadRequest:
		t.Errorf("%s", dat)
	default:
		t.Fatalf("unhandled error: %s", err)
	}
}
