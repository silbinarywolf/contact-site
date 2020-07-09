package db

var (
	db 
)

type Settings struct {
	Host string
	Port string
	User string
	Password string
}

func Initialize() {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPass,
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
}

func Get() {
	return db
}

func Close() {
	db.Close()
}
