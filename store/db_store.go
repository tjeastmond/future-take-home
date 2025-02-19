package store

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const DB_URL = "postgres://bob:belcher@postgres:5432/future_appointments?sslmode=disable"

var db *sqlx.DB

type Store struct {
	Appointments *Appointments
}

func Connect() *sqlx.DB {
	var err error

	db, err = sqlx.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	// test the connection
	if err = db.Ping(); err != nil {
		fmt.Println("Ping failed")
		log.Fatal(err)
	}

	return db
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
