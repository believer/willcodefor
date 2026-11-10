package data

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/dotsql"
	"github.com/swithek/dotsqlx"

	_ "github.com/lib/pq"
)

var (
	DB      *sqlx.DB
	UmamiDB *sqlx.DB
	Dot     *dotsqlx.DotSqlx
)

func InitDB() error {
	connectionString := os.Getenv("DATABASE_URL")

	db := sqlx.MustConnect("postgres", connectionString)
	err := db.Ping()

	if err != nil {
		return err
	} else {
		log.Println("Connected to database")
	}

	// Optional, view counts only include frozen post_view data without it
	if umamiConnectionString := os.Getenv("UMAMI_DATABASE_URL"); umamiConnectionString != "" {
		UmamiDB, err = sqlx.Connect("postgres", umamiConnectionString)

		if err != nil {
			return err
		}

		log.Println("Connected to Umami database")
	}

	indexQueries, err := dotsql.LoadFromFile("./data/indexQueries.sql")

	if err != nil {
		return err
	}

	postQueries, err := dotsql.LoadFromFile("./data/postQueries.sql")

	if err != nil {
		return err
	}

	bookQueries, err := dotsql.LoadFromFile("./data/bookQueries.sql")

	if err != nil {
		return err
	}

	umamiQueries, err := dotsql.LoadFromFile("./data/umamiQueries.sql")

	if err != nil {
		return err
	}

	dot := dotsql.Merge(
		indexQueries,
		postQueries,
		bookQueries,
		umamiQueries,
	)
	dotx := dotsqlx.Wrap(dot)

	// Set the global DBClient variable to the db connection
	DB = db
	Dot = dotx

	return nil

}
