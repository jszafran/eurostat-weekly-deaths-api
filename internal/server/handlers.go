package server

import (
	"database/sql"
	"log"
	"net/http"
	"weekly_deaths/internal/db"
	"weekly_deaths/internal/queries"
)

func CountriesHandler(w http.ResponseWriter, r *http.Request) {
	var (
		week    int
		year    int
		deaths  sql.NullInt64
		age     string
		sex     string
		country string
	)
	countryParam := r.URL.Query().Get("country")
	if countryParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("country url param required"))
		return
	}

	genderParam := r.URL.Query().Get("gender")
	if genderParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("gender url param required"))
		return
	}

	ageParam := r.URL.Query().Get("age")
	if ageParam != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("age url param required"))
		return
	}

	db, err := db.GetDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare(queries.WEEKLY_DEATHS_FOR_COUNTRY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var weeklyDeaths []WeeklyDeaths
	rows, err := stmt.Query(countryParam, genderParam, ageParam)
	for rows.Next() {
		err := rows.Scan(&week, &year, &deaths, &age, &sex, &country)
		if err != nil {
			log.Fatal(err)
		}
		weeklyDeaths = append(weeklyDeaths, WeeklyDeaths{
			Week:         week,
			Year:         year,
			WeeklyDeaths: deaths,
			Age:          age,
			Sex:          sex,
			Country:      country,
		})
	}
	defer rows.Close()
	log.Printf("%+v\n", weeklyDeaths)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"hello": "world"}`))
}
