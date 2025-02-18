package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const DB_URL = "postgres://bob:belcher@postgres:5432/future_appointments?sslmode=disable"

var db *sql.DB

func Connect() {
	var err error

	db, err = sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	// test the connection
	if err = db.Ping(); err != nil {
		fmt.Println("Ping failed")
		log.Fatal(err)
	}
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func GetDB() *sql.DB {
	return db
}
