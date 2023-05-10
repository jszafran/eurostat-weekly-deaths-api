package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"weekly_deaths/internal/db"
)

func CountriesHandler(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		WriteJSON(http.StatusBadRequest, w, map[string]string{"error": "country url param required"})
		return
	}

	gender := r.URL.Query().Get("gender")
	if gender == "" {
		WriteJSON(http.StatusBadRequest, w, map[string]string{"error": "gender url param required"})
		return
	}

	age := r.URL.Query().Get("age")
	if age == "" {
		WriteJSON(http.StatusBadRequest, w, map[string]string{"error": "age url param required"})
		return
	}

	yearFromStr := r.URL.Query().Get("year_from")
	if yearFromStr == "" {
		WriteJSON(http.StatusBadRequest, w, map[string]string{"error": "year_from url param required"})
		return
	}
	yearFrom, err := strconv.Atoi(yearFromStr)
	if err != nil {
		WriteJSON(http.StatusBadRequest, w, map[string]string{
			"error": fmt.Sprintf("value %s cannot be converted to int", yearFromStr),
		})
		return
	}

	yearToStr := r.URL.Query().Get("year_to")
	if yearToStr == "" {
		WriteJSON(http.StatusBadRequest, w, map[string]string{"error": "year_to url param required"})
		return
	}
	yearTo, err := strconv.Atoi(yearToStr)
	if err != nil {
		WriteJSON(http.StatusBadRequest, w, map[string]string{
			"error": fmt.Sprintf("value %s cannot be converted to int", yearToStr),
		})
		return
	}

	database, err := db.GetDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weeklyDeaths, err := db.GetCountryData(database, country, gender, age, yearFrom, yearTo)
	if err != nil {
		log.Println(err)
		WriteJSON(http.StatusInternalServerError, w, map[string]string{"error": "internal server error"})
		return
	}

	data := db.WeeklyDeathsResponse{Data: weeklyDeaths}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
