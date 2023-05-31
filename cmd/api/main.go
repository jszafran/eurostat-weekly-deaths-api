package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"weekly_deaths/internal/eurostat"
)

// DefaultPort defines a default port that the server will be started on.
const DefaultPort = 8080

type application struct {
	port             int
	dataProvider     *eurostat.DataProvider
	dataDownloadedAt time.Time
}

func main() {
	var err error
	var port int

	flag.IntVar(&port, "port", DefaultPort, "port to start server on")
	flag.Parse()

	dp, err := eurostat.NewDataProvider(eurostat.LiveEurostatDataSource{})
	if err != nil {
		log.Fatal(err)
	}

	app := application{
		port:             port,
		dataProvider:     dp,
		dataDownloadedAt: time.Now(),
	}

	router := app.routes()
	log.Printf("Starting the server on :%d port\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
