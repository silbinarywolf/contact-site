package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	config Config
)

type Config struct {
	Database struct {
		Host     string `json:"host,omitempty"`
		Port     int    `json:"port,omitempty"`
		User     string `json:"user,omitempty"`
		Password string `json:"password,omitempty"`
	}
}

// Get will return a copy of the current application configuration.
// MustLoad must be called before this is called.
func Get() Config {
	return config
}

func MustLoad() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	dat, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Config load error: %s\n", err)
	}
	// NOTE(Jae): 2020-07-11
	// I considered using *.toml as I prefer that format over JSON.
	// But in the interest of keeping external dependencies down and things simple,
	// I decided to just use *.json.
	var newConfig Config
	if err := json.Unmarshal(dat, &newConfig); err != nil {
		log.Fatalf("Config parse error: %s\n", err)
	}
	// Print all the config errors we get at once, rather than one at a time to make resolving
	// potential configuration mistakes nicer.
	shouldEarlyExit := false
	if newConfig.Database.User == "" {
		log.Printf("\"user\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Password == "" {
		log.Printf("\"password\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Host == "" {
		log.Printf("\"host\" JSON key for environment variable cannot be empty.")
		shouldEarlyExit = true
	}
	if newConfig.Database.Port == 0 {
		log.Printf("\"port\" JSON key for environment variable cannot be empty or set to 0.")
		shouldEarlyExit = true
	}
	if shouldEarlyExit {
		os.Exit(1)
	}
	// Set config on success
	config = newConfig
}
