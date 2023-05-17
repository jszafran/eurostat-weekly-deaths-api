package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"weekly_deaths/internal/db"
	"weekly_deaths/internal/eurostat"
)

func WeeklyDeathsHandler(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		WriteJSONError(http.StatusBadRequest, w, "country url param required")
		return
	}

	gender := r.URL.Query().Get("gender")
	if gender == "" {
		WriteJSONError(http.StatusBadRequest, w, "gender url param required")
		return
	}

	age := r.URL.Query().Get("age")
	if age == "" {
		WriteJSONError(http.StatusBadRequest, w, "age url param required")
		return
	}

	yearFromStr := r.URL.Query().Get("year_from")
	if yearFromStr == "" {
		WriteJSONError(http.StatusBadRequest, w, "year_from url param required")
		return
	}
	yearFrom, err := strconv.Atoi(yearFromStr)
	if err != nil {
		WriteJSONError(http.StatusBadRequest, w, fmt.Sprintf("value %s cannot be converted to int", yearFromStr))
		return
	}

	yearToStr := r.URL.Query().Get("year_to")
	if yearToStr == "" {
		WriteJSONError(http.StatusBadRequest, w, "year_to url param required")
		return
	}
	yearTo, err := strconv.Atoi(yearToStr)
	if err != nil {
		WriteJSONError(http.StatusBadRequest, w, fmt.Sprintf("value %s cannot be converted to int", yearToStr))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weeklyDeaths, err := db.GetCountryData(country, gender, age, yearFrom, yearTo)
	if err != nil {
		WriteJSONError(http.StatusInternalServerError, w, "internal server error")
		return
	}

	data := db.WeeklyDeathsResponse{Data: weeklyDeaths}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func LabelsHandler(w http.ResponseWriter, r *http.Request) {
	data := eurostat.GetLabels()
	WriteJSON(http.StatusOK, w, map[string][]eurostat.MetadataLabel{"data": data})
}
