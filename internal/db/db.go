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

// Get will get an active database connection.
//
// The general pattern in this codebase is to just use this to get the current database
// connection rather than passing around a *sql.DB pointer everywhere. I'm betting that
// I'm not going to need multiple database connections/drivers in this project.
//
// I have a hunch that this is a bad idea and that leveraging the "context.Context" object
// to get the current DB object would allow for more code-reuse but right now from the code I have
// its unclear as to what practical benefits that would give me and I'm out of time.
//
// Safe for concurrent use.
// (As per Golang docs, *sql.DB is safe for concurrent use)
func Get() *sql.DB {
	return db
}

func MustConnect(settings Settings) {
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
	//
	// Simply calling "Open" won't tell us if it connected successfully, at least for the Postgres
	// driver I'm using. The reason for the retry logic is to improve support when running in a
	// Docker environment.
	for i := 0; i < maxDBRetries; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Printf("Database connection attempt #%d: %v\n", i, err)
		if i == maxDBRetries-1 {
			// Not panicing here as its not a developer-fault, this kind of error
			// is likely to be a user/config error and so we don't need the callstack.
			log.Println("Unable to connect to database. Stopping app.")
			os.Exit(1)
		}
		time.Sleep(timeBetweenDBRetries)
	}
}

func MustClose() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
		db = nil
	}
	return nil
}
