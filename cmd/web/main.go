package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"weekly_deaths/internal/server"
)

func main() {
	r := chi.NewRouter()
	r.Get("/api/countries", server.CountriesHandler)
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
