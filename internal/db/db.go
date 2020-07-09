package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	db *sql.DB
)

type Settings struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

func Connect(settings Settings) {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		settings.Host,
		settings.Port,
		settings.User,
		settings.Password,
	))
	if err != nil {
		panic(err)
	}

	// Test connection to the database
	for i := 0; i < 5; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Printf("Database connection attempt #%d: %v\n", i, err)
		if i == 4 {
			log.Println("Unable to connect to database. Stopping app.")
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}

	// Select database
	/*if _, err := db.Query("SELECT DATABASE " + settings.DatabaseName + ";"); err != nil {
		log.Printf("Unable to select database: %s\n", err)
		os.Exit(1)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			// Do nothing if "duplicate_database" error
			// it's already been created
		} else {
			panic(err)
		}
	}*/
	//log.Println("Database connection successful")
}

func Get() *sql.DB {
	return db
}

func Close() {
	db.Close()
}
