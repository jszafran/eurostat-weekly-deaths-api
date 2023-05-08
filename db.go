package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// GetDB prepares a connection to SQLite database.
func GetDB() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "./eurostat.db")
	if err != nil {
		return db, err
	}

	return db, err
}

// RecreateDB drops and creates a weekly_deaths table.
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

func InsertWeeklyDeathsData(records []WeeklyDeathsRecord, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO weekly_deaths (week, year, deaths, age, sex, country)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, r := range records {
		_, err = stmt.Exec(
			r.Week,
			r.Year,
			r.Deaths,
			r.Age,
			r.Sex,
			r.Country,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
