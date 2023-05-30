// Package server contains all code related to creating
// HTTP server and serving Eurostat data through a simple API.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"weekly_deaths/internal/eurostat"
)

// Commit variable stores an information about
// current commit that the application was built from.
var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}()

// WeeklyDeathsHandler is a HTTP handler func exposing
// Eurostat weekly deaths data for given parameters:
// - country
// - gender
// - age
// - year_from
// - year_to
// All parameters are required and should be passed as query params.
func WeeklyDeathsHandler(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		writeJSONError(http.StatusBadRequest, w, "country url param required")
		return
	}

	gender := r.URL.Query().Get("gender")
	if gender == "" {
		writeJSONError(http.StatusBadRequest, w, "gender url param required")
		return
	}

	age := r.URL.Query().Get("age")
	if age == "" {
		writeJSONError(http.StatusBadRequest, w, "age url param required")
		return
	}

	yearFromStr := r.URL.Query().Get("year_from")
	if yearFromStr == "" {
		writeJSONError(http.StatusBadRequest, w, "year_from url param required")
		return
	}
	yearFrom, err := strconv.Atoi(yearFromStr)
	if err != nil {
		writeJSONError(http.StatusBadRequest, w, fmt.Sprintf("value %s cannot be converted to int", yearFromStr))
		return
	}

	yearToStr := r.URL.Query().Get("year_to")
	if yearToStr == "" {
		writeJSONError(http.StatusBadRequest, w, "year_to url param required")
		return
	}
	yearTo, err := strconv.Atoi(yearToStr)
	if err != nil {
		writeJSONError(http.StatusBadRequest, w, fmt.Sprintf("value %s cannot be converted to int", yearToStr))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weeklyDeaths, err := eurostat.EurostatDataProvider.GetWeeklyDeaths(
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
func LabelsHandler(w http.ResponseWriter, r *http.Request) {
	data := eurostat.GetLabels()
	writeJSON(http.StatusOK, w, map[string][]eurostat.MetadataLabel{"data": data})
}

// InfoHandler is an HTTP handler returnign metadata about the application:
// - the commit from which currently running istance was built
// - timestamp indicating when the data was downloaded from Eurostat
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(http.StatusOK, w, map[string]any{
		"commit_hash":                      Commit,
		"data_downloaded_at_utc_timestamp": eurostat.DataDownloadedAt,
	})
}
