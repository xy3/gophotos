package photos

import (
	"database/sql"
	"io/ioutil"
	"log"
)

var (
	DB *sql.DB
)

func DbSetup(db *sql.DB) error {
	log.Print("Setting up database...")
	query, err := ioutil.ReadFile("schema/schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err = db.Exec(string(query)); err != nil {
		panic(err)
	}
	log.Print("Set up the database successfully")
	return nil
}

func DbConnect() (err error) {
	DB, err = sql.Open("sqlite3", "photos.sqlite")
	return
}
