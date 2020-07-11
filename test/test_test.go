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
)

func TestMain(m *testing.M) {
	// Start application
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get current dir: %s", err))
	}
	if err := os.Chdir(filepath.Join(dir, "..")); err != nil {
		panic(fmt.Sprintf("failed to change dir: %s", err))
	}
	go app.Start()

	// Run tests
	m.Run()
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
	//panic(string(dat))
}

func TestPostForm(t *testing.T) {
	//type TestData struct {
	//	In  string
	//	Out bool
	//}
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

	//t.Logf("%s", dat)
	//panic(string(dat))
}
