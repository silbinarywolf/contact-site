package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	// It might be more appropriate for these constants to be configurable instead,
	// but we can always decide to do that later. For now, this is probably good enough.
	maxDBRetries         = 5
	timeBetweenDBRetries = 2 * time.Second
)

var (
	db *sql.DB
)

type Settings struct {
	Host     string
	Port     int
	User     string
	Password string
	// Opted to not implement for time reasons. We just use Postgres's default database.
	// DatabaseName string
}

// Get will get an active database connection
//
// Safe for concurrent use.
func Get() *sql.DB {
	return db
}

// Connect will connect to a postgres database
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
	for i := 0; i < maxDBRetries; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Printf("Database connection attempt #%d: %v\n", i, err)
		if i == maxDBRetries-1 {
			log.Println("Unable to connect to database. Stopping app.")
			os.Exit(1)
		}
		time.Sleep(timeBetweenDBRetries)
	}
}

func Close() {
	db.Close()
}
