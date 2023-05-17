package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"weekly_deaths/internal/db"
	"weekly_deaths/internal/server"
)

const port = ":3000"

func main() {
	var err error

	db.DB, err = db.GetDB()
	r := chi.NewRouter()

	r.Get("/api/weekly_deaths", server.WeeklyDeathsHandler)
	r.Get("/api/labels", server.LabelsHandler)

	log.Printf("Starting the server on %s port\n", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("starting server: %s\n", err.Error())
	}
}
