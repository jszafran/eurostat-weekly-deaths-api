package main

import (
	"log"
	"time"

	"weekly_deaths/internal/db"
	"weekly_deaths/internal/eurostat"
)

func main() {
	var err error

	t1 := time.Now()

	db.DB, err = db.GetDB()
	if err != nil {
		log.Fatalf("creating a database: %s\n", err)
	}

	err = db.RecreateWeeklyDeathsTable()
	if err != nil {
		log.Fatalf("recreating weekly_deaths table: %s\n", err)
	}

	data, err := eurostat.ReadData()
	if err != nil {
		log.Fatalf("reading data from eurostat: %s\n", err)
	}

	recs, err := eurostat.ParseData(data)
	log.Printf("Parsed %d records.\n", len(recs))
	log.Println("Starting inserting data into db.")

	err = db.InsertWeeklyDeathsData(recs)
	if err != nil {
		log.Fatal(err)
	}

	var recordsInserted int
	err = db.DB.QueryRow("select count(*) from weekly_deaths").Scan(&recordsInserted)
	if err != nil {
		log.Fatalf("reading inserted records count from db: %s\n", err)
	}

	dbix := "idx_weekly_deaths"
	log.Printf("Recreating %s index.\n", dbix)

	err = db.DropObject(dbix, db.DROP_INDEX_SQL)
	if err != nil {
		log.Fatalf("dropping %s index: %s", dbix, err)
	}

	err = db.CreateObject(dbix, db.CREATE_INDEX_SQL)
	if err != nil {
		log.Fatalf("creating %s index: %s", dbix, err)
	}

	if recordsInserted != len(recs) {
		log.Fatalf("%d records were not inserted into weekly_deaths table.", len(recs)-recordsInserted)
	}

	log.Println("Cleaning incorrect data (week > 53).")
	_, err = db.DB.Exec(db.DELETE_INCORRECT_WEEKS_DATA_SQL)
	if err != nil {
		log.Fatalf("deleting incorrect week data: %s", err)
	}

	err = eurostat.ValidateLabels(recs)
	if err != nil {
		log.Fatal(err)
	}

	timeElapsed := time.Since(t1)

	log.Printf("%d records inserted. Time elapsed: %s.\n", recordsInserted, timeElapsed)
}
