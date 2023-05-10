package db

import (
	"database/sql"
	"fmt"
	"log"
	"weekly_deaths/internal/eurostat"
	"weekly_deaths/internal/queries"

	_ "github.com/mattn/go-sqlite3"
)

type WeeklyDeathsResponse struct {
	Data []WeeklyDeaths `json:"data"`
}

type WeeklyDeaths struct {
	Year         int           `json:"year"`
	Week         int           `json:"week"`
	WeeklyDeaths sql.NullInt64 `json:"weekly_deaths"`
	Age          string        `json:"age"`
	Sex          string        `json:"sex"`
	Country      string        `json:"country"`
}

// GetDB prepares a connection to SQLite database.
func GetDB() (*sql.DB, error) {
	var db *sql.DB
	db, err := sql.Open("sqlite3", "../../eurostat.db")
	if err != nil {
		return db, err
	}

	return db, err
}

func RecreateTable(table string, ddlQuery string, db *sql.DB) error {
	log.Printf("Recreating %s table.\n", table)
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	if err != nil {
		return err
	}

	_, err = db.Exec(ddlQuery)
	if err != nil {
		return err
	}

	return nil
}

// RecreateDB drops and creates a weekly_deaths, countries, ages, genders tables.
func RecreateDB(db *sql.DB) error {
	err := RecreateTable("weekly_deaths", queries.CREATE_WEEKLY_DEATHS_SQL, db)
	if err != nil {
		return err
	}

	err = RecreateTable("countries", queries.CREATE_COUNTRIES_SQL, db)
	if err != nil {
		return err
	}

	err = RecreateTable("ages", queries.CREATE_AGES_SQL, db)
	if err != nil {
		return err
	}

	err = RecreateTable("genders", queries.CREATE_GENDERS_SQL, db)
	if err != nil {
		return err
	}

	return nil
}

func PopulateMetadataTables(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO countries SELECT DISTINCT country FROM weekly_deaths")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO ages SELECT DISTINCT age FROM weekly_deaths")
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO genders SELECT DISTINCT sex FROM weekly_deaths")
	if err != nil {
		return err
	}

	return nil
}

func InsertWeeklyDeathsData(records []eurostat.WeeklyDeathsRecord, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO weekly_deaths (week, year, deaths, age, sex, country)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	for _, r := range records {
		_, err = stmt.Exec(
			r.Week,
			r.Year,
			r.Deaths,
			r.Age,
			r.Sex,
			r.Country,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetCountryData(db *sql.DB, countryParam string, sexParam string, ageParam string) ([]WeeklyDeaths, error) {
	var (
		week    int
		year    int
		deaths  sql.NullInt64
		age     string
		sex     string
		country string
		results []WeeklyDeaths
	)

	stmt, err := db.Prepare(queries.WEEKLY_DEATHS_FOR_COUNTRY)
	if err != nil {
		return results, err
	}

	rows, err := stmt.Query(countryParam, sexParam, ageParam)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&week, &year, &deaths, &age, &sex, &country)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, WeeklyDeaths{
			Week:         week,
			Year:         year,
			WeeklyDeaths: deaths,
			Age:          age,
			Sex:          sex,
			Country:      country,
		})
	}

	return results, err
}
