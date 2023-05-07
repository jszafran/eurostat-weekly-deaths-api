package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./eurostat.db")
	if err != nil {
		return db, err
	}

	return db, err
}

func RecreateDB(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS weekly_deaths")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE weekly_deaths (
			week INTEGER NOT NULL,
			year INTEGER NOT NULL,
			deaths INTEGER,
			age STRING,
			sex STRING,
			country STRING
		) 
	`)
	if err != nil {
		return err
	}

	return nil
}
