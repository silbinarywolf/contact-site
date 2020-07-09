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
	var newConfig Config
	if err := json.Unmarshal(dat, &newConfig); err != nil {
		log.Fatalf("Config parse error: %s\n", err)
	}
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

func Get() Config {
	return config
}
