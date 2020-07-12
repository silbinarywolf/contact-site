package config

import (
	"encoding/json"
	"log"
	"os"
)

const (
	configBasename = "config.json"
)

var (
	config Config
)

// Config structure that maps to a configuration file.
type Config struct {
	Web struct {
		Port int `json:"port,omitempty"`
	} `json:"web,omitempty"`
	Database struct {
		Host     string `json:"host,omitempty"`
		Port     int    `json:"port,omitempty"`
		User     string `json:"user,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"database,omitempty"`
}

// Get will return a copy of the current application configuration.
// MustLoad must be called before this is called.
func Get() Config {
	return config
}

// Exists will check if the config file exists or not.
// This was implemented for use by integration tests.
func Exists() bool {
	if _, err := os.Stat(configBasename); os.IsNotExist(err) {
		return false
	}
	return true
}

// MustLoad will try to load the applications config.json file.
// Panics if an error occurs.
//
// This could return an error and be handled that way, but I opted to not
// put the work as of yet as there's no clear benefit in doing so. If I needed this
// to be easier to work with in tests or something, it'd probably be worth it, in which
// case there would be a Load function implemented, and MustLoad would just wrap it and panic
// if error is not nil.
func MustLoad() {
	file, err := os.Open(configBasename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// NOTE(Jae): 2020-07-11
	// I considered using *.toml as I prefer that format over JSON.
	// But in the interest of keeping external dependencies down and things simple,
	// I decided to just use *.json.
	var newConfig Config
	{
		decoder := json.NewDecoder(file)
		// I typo things all the time, so I want to know if my configuration file is trying to use
		// a key/field that doesn't exist as soon as possible.
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&newConfig); err != nil {
			log.Fatalf("Config parse error: %s\n", err)
		}
	}
	// Print all the config errors we get at once, rather than one at a time to make resolving
	// potential configuration mistakes nicer.
	//
	// Might be a good idea to change this in the future to leverage reflection for checking each key
	// but this is good enough for now.
	shouldEarlyExit := false
	if newConfig.Web.Port == 0 {
		log.Printf("\"web.port\" JSON key for environment variable cannot be empty or set to 0.")
		shouldEarlyExit = true
	}
	if newConfig.Database.User == "" {
		log.Printf("\"database.user\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Password == "" {
		log.Printf("\"database.password\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Host == "" {
		log.Printf("\"database.host\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Port == 0 {
		log.Printf("\"database.port\" JSON key for environment variable cannot be empty or set to 0.")
		shouldEarlyExit = true
	}
	if shouldEarlyExit {
		os.Exit(1)
	}
	// Set config on success
	config = newConfig
}
