package main

import (
	"flag"

	"github.com/silbinarywolf/contact-site/internal/app"
)

func main() {
	// Parse flags here to avoid conflicts with test flags
	// not being parsed
	flag.Parse()

	// Put the application in its own package, this will give us the ability to run
	// the entire application as an integration test later and possibly extend the
	// Start method and pass in a "test mode" specific flag
	app.Start()
}
