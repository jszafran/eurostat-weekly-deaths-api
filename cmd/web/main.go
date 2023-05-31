// Package downloads data from Eurostat and
// starts HTTP server to expose this data
// (server listens on 8080 port by default).
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"weekly_deaths/internal/eurostat"
	"weekly_deaths/internal/server"
)

// DefaultPort defines a default port that the server will be started on.
const DefaultPort = 8080

func main() {
	var err error
	var port int

	flag.IntVar(&port, "port", DefaultPort, "port to start server on")
	flag.Parse()

	eurostat.EurostatDataProvider, err = eurostat.NewDataProvider(eurostat.LiveEurostatDataSource{})
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Get("/api/weekly_deaths", server.WeeklyDeathsHandler)
	r.Get("/api/labels", server.LabelsHandler)
	r.Get("/api/info", server.InfoHandler)

	log.Printf("Starting the server on :%d port\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
