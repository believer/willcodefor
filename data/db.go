package data

import (
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		panic(err)
	}

	// Check if we can connect to the DB
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	DB = db
}
