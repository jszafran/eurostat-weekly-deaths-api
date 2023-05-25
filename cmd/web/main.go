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

func main() {
	var err error
	var port int

	flag.IntVar(&port, "port", 3000, "port to start server on")
	flag.Parse()

	eurostat.EurostatDataProvider, err = eurostat.NewDataProvider(eurostat.LiveEurostatDataSource{})
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Get("/api/weekly_deaths", server.WeeklyDeathsHandler)
	r.Get("/api/labels", server.LabelsHandler)

	log.Printf("Starting the server on :%d port\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
