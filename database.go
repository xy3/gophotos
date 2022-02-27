package photos

import (
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	DB *sql.DB
	//go:embed schema/schema.sql
	setupSQL string
)

func DbSetup(db *sql.DB) error {
	log.Println("Setting up database...")
	if _, err := db.Exec(setupSQL); err != nil {
		panic(err)
	}
	log.Println("Set up the database successfully")
	return nil
}

func DbConnect() (err error) {
	DB, err = sql.Open("sqlite3", Config.SqliteDSN)
	DB.SetMaxOpenConns(1)
	return
}
