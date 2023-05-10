package server

import (
	"encoding/json"
	"log"
	"net/http"
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

	database, err := db.GetDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	weeklyDeaths, err := db.GetCountryData(database, country, gender, age)
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
