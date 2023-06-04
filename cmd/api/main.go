package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"weekly_deaths/internal/eurostat"
)

// DefaultPort defines a default port that the server will be started on.
const DefaultPort = 8080

func main() {
	var err error
	var port int

	flag.IntVar(&port, "port", DefaultPort, "port to start server on")
	flag.Parse()

	if err != nil {
		log.Fatal(err)
	}

	snapshot, err := eurostat.DataSnapshotFromPath("../../snapshots/20230604T114323.tsv.gz")
	if err != nil {
		log.Fatal(err)
	}
	db := eurostat.DBFromSnapshot(snapshot)
	if err != nil {
		log.Fatal(err)
	}

	app := application{
		db: db,
	}

	router := app.routes()
	log.Printf("Starting the server on :%d port\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
