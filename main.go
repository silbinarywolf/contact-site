package main

import (
	"github.com/silbinarywolf/contact-site/internal/app"
)

func main() {
	// Put the application in its own package, this will give us the ability to run
	// the entire application as an integration test later and possibly extend the
	// Start method to pass in a "test mode" specific flag
	app.Start()
}
