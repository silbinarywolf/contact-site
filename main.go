package main

import (
	"flag"

	"github.com/silbinarywolf/contact-site/internal/app"
)

func main() {
	// Parse flags here to avoid conflicts with test flags
	flag.Parse()

	// Put the application in its own package, this will give us the ability to run
	// the entire application as an integration test
	//
	// We seperate the initialization and startup of the server for test purposes
	app.MustInitialize()
	defer app.MustClose()

	app.MustStart()
}
