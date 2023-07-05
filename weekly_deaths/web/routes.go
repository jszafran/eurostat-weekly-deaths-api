package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"weekly_deaths/eurostat"

	"github.com/go-chi/chi/v5"
)

const ParamRequiredUserMessage = "This query url parameter is required."
const FailedConversionToIntMessage = "Provided value cannot be converted to integer."

type Application struct {
	Db   *eurostat.InMemoryDB
	Auth struct {
		Username string
		Password string
	}
}

func (app *Application) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/api/weekly_deaths", app.WeeklyDeathsHandler)
	router.Get("/api/labels", app.LabelsHandler)
	router.Get("/api/info", app.InfoHandler)
	router.Post("/api/update_data", app.basicAuth(app.UpdateDataHandler))

	return router
}

// WeeklyDeathsHandler is a HTTP handler func exposing
// Eurostat weekly deaths data for given parameters:
// - country
// - gender
// - age
// - year_from
// - year_to
// All parameters are required and should be passed as query params.
func (app *Application) WeeklyDeathsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		yearFrom int
		yearTo   int
		err      error
	)
	errors := make(map[string]string)

	country := r.URL.Query().Get("country")
	if country == "" {
		errors["country"] = ParamRequiredUserMessage
	}

	gender := r.URL.Query().Get("gender")
	if gender == "" {
		errors["gender"] = ParamRequiredUserMessage
	}

	age := r.URL.Query().Get("age")
	if age == "" {
		errors["age"] = ParamRequiredUserMessage
	}

	yearFromStr := r.URL.Query().Get("year_from")
	if yearFromStr == "" {
		errors["year_from"] = ParamRequiredUserMessage
	} else {
		yearFrom, err = strconv.Atoi(yearFromStr)
		if err != nil {
			errors["year_from"] = FailedConversionToIntMessage
		}
	}

	yearToStr := r.URL.Query().Get("year_to")
	if yearToStr == "" {
		errors["year_to"] = ParamRequiredUserMessage
	} else {
		yearTo, err = strconv.Atoi(yearToStr)
		if err != nil {
			errors["year_to"] = FailedConversionToIntMessage
		}
	}

	if len(errors) > 0 {
		writeJSONError(http.StatusBadRequest, w, errors)
		return
	}

	weeklyDeaths, err := app.Db.GetWeeklyDeaths(
		country,
		age,
		gender,
		yearFrom,
		yearTo,
	)
	if err != nil {
		writeJSONError(http.StatusInternalServerError, w, "internal server error")
		return
	}

	data := WeeklyDeathsResponse{Gender: gender, Age: age, Country: country, WeeklyDeaths: weeklyDeaths}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// LabelsHandler is an HTTP handler returning labels translation
// for countries, genders and age groups used in weekly deaths dataset.
func (app *Application) LabelsHandler(w http.ResponseWriter, r *http.Request) {
	data := GetLabels()
	writeJSON(http.StatusOK, w, map[string][]MetadataLabel{"data": data})
}

// InfoHandler is an HTTP handler returning metadata about the application:
// - the commit from which currently running instance was built
// - timestamp indicating when the data was downloaded from Eurostat
func (app *Application) InfoHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(http.StatusOK, w, InfoResponse{
		CommitHash:       os.Getenv("COMMIT"),
		DataDownloadedAt: app.Db.Timestamp(),
	})
}

func (app *Application) UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for data update.")
	snapshot, err := eurostat.DataSnapshotFromEurostat()
	if err != nil {
		log.Printf("Data update failed: %s\n", err)
		writeJSONError(http.StatusInternalServerError, w, fmt.Sprintf("Fetching data from Eurostat failed: %s", err))
	}

	app.Db.LoadSnapshot(snapshot)
	log.Println("Data update succeeded.")
	msg := fmt.Sprintf("Successfully loaded snapshot for %s.", snapshot.Timestamp)
	writeJSON(http.StatusOK, w, map[string]string{"message": msg})
}
