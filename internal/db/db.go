package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"weekly_deaths/internal/eurostat"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// JsonNullInt64 marshals just the integer value (instead of the Valid/NullInt64 wrapper).
// Credits: https://stackoverflow.com/questions/33072172/how-can-i-work-with-sql-null-values-and-json-in-a-good-way/33072822#33072822
type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type WeeklyDeathsResponse struct {
	Data []WeeklyDeaths `json:"data"`
}

type WeeklyDeaths struct {
	Year         int           `json:"year"`
	Week         int           `json:"week"`
	WeeklyDeaths JsonNullInt64 `json:"weekly_deaths"`
	Age          string        `json:"age"`
	Gender       string        `json:"gender"`
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

// Recreate table drops and creates back given table.
func RecreateTable(table string, ddlQuery string) error {
	log.Printf("Recreating %s table.\n", table)
	_, err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	if err != nil {
		return fmt.Errorf("dropping %s table: %w\n", table, err)
	}

	_, err = DB.Exec(ddlQuery)
	if err != nil {
		return fmt.Errorf("creating %s table: %w\n", table, err)
	}

	return nil
}

// RecreateDB drops and creates a weekly_deaths, countries, ages, genders tables.
func RecreateDB() error {
	err := RecreateTable("weekly_deaths", CREATE_WEEKLY_DEATHS_SQL)
	if err != nil {
		return fmt.Errorf("recreating weekly_deaths table: %w\n", err)
	}

	return nil
}

func InsertWeeklyDeathsData(records []eurostat.WeeklyDeathsRecord) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("inserting weekly deaths data - beginning transaction: %w\n", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO weekly_deaths (week, year, deaths, age, gender, country)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("preparing insert statement for weekly deaths: %w", err)
	}

	defer stmt.Close()
	for _, r := range records {
		_, err = stmt.Exec(
			r.Week,
			r.Year,
			r.Deaths,
			r.Age,
			r.Gender,
			r.Country,
		)
		if err != nil {
			return fmt.Errorf(
				"executing insert for (%d, %d, %v, %s, %s, %s): %w\n",
				r.Week,
				r.Year,
				r.Deaths,
				r.Age,
				r.Gender,
				r.Country,
				err,
			)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transction for weekly deaths insert: %w\n", err)
	}

	return nil
}

func GetCountryData(
	countryParam string,
	genderParam string,
	ageParam string,
	yearFrom int,
	yearTo int,
) ([]WeeklyDeaths, error) {
	var (
		week    int
		year    int
		deaths  JsonNullInt64
		age     string
		gender  string
		country string
		results []WeeklyDeaths
	)

	stmt, err := DB.Prepare(WEEKLY_DEATHS_FOR_COUNTRY)
	if err != nil {
		return results, err
	}

	rows, err := stmt.Query(countryParam, genderParam, ageParam, yearFrom, yearTo)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&week, &year, &deaths, &age, &gender, &country)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, WeeklyDeaths{
			Week:         week,
			Year:         year,
			WeeklyDeaths: deaths,
			Age:          age,
			Gender:       gender,
			Country:      country,
		})
	}

	return results, nil
}
