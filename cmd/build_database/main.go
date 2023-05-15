package main

import (
	"log"
	"time"

	"weekly_deaths/internal/db"
	"weekly_deaths/internal/eurostat"
)

func main() {
	t1 := time.Now()
	database, err := db.GetDB()
	if err != nil {
		log.Fatalf("creating a database: %s\n", err)
	}

	err = db.RecreateDB(database)
	if err != nil {
		log.Fatalf("recreating database tables: %s\n", err)
	}

	data, err := eurostat.ReadData()
	if err != nil {
		log.Fatalf("reading data from eurostat: %s\n", err)
	}

	recs, err := eurostat.ParseData(data)
	log.Printf("Parsed %d records.\n", len(recs))
	log.Println("Starting inserting data into db.")

	err = db.InsertWeeklyDeathsData(recs, database)
	if err != nil {
		log.Fatal(err)
	}

	var recordsInserted int
	err = database.QueryRow("select count(*) from weekly_deaths").Scan(&recordsInserted)
	if err != nil {
		log.Fatalf("reading inserted records count from db: %s\n", err)
	}

	timeElapsed := time.Since(t1)

	log.Printf("%d records inserted. Time elapsed: %s.\n", recordsInserted, timeElapsed)

	if recordsInserted != len(recs) {
		log.Fatalf("%d records were not inserted into weekly_deaths table.", len(recs)-recordsInserted)
	}

}
