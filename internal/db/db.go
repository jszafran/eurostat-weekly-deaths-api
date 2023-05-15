package db

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"weekly_deaths/internal/eurostat"

	_ "github.com/mattn/go-sqlite3"
)

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

type MetadataLabel struct {
	Value string
	Label string
	Order int
}

type MetadataLabelFromDB struct {
	Value     string
	Label     string
	Order     int
	LabelType string
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
func RecreateTable(table string, ddlQuery string, db *sql.DB) error {
	log.Printf("Recreating %s table.\n", table)
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	if err != nil {
		return fmt.Errorf("dropping %s table: %w\n", table, err)
	}

	_, err = db.Exec(ddlQuery)
	if err != nil {
		return fmt.Errorf("creating %s table: %w\n", table, err)
	}

	return nil
}

// RecreateDB drops and creates a weekly_deaths, countries, ages, genders tables.
func RecreateDB(db *sql.DB) error {
	err := RecreateTable("weekly_deaths", CREATE_WEEKLY_DEATHS_SQL, db)
	if err != nil {
		return fmt.Errorf("recreating weekly_deaths table: %w\n", err)
	}

	err = RecreateTable("countries", CREATE_COUNTRIES_SQL, db)
	if err != nil {
		return fmt.Errorf("recreating countries table: %w\n", err)
	}

	err = RecreateTable("ages", CREATE_AGES_SQL, db)
	if err != nil {
		return fmt.Errorf("recreating ages table: %w\n", err)
	}

	err = RecreateTable("genders", CREATE_GENDERS_SQL, db)
	if err != nil {
		return fmt.Errorf("recreating genders table: %w\n", err)
	}

	return nil
}

func LoadMetadataTable(db *sql.DB, csvPath string, table string) error {
	labels, err := parseLabelFile(csvPath)
	if err != nil {
		return fmt.Errorf("parsing labels from csv: %w\n", err)
	}

	if len(labels) == 0 {
		log.Fatalf("Labels CSV: no records parsed for %s table.\n", table)
	}

	err = InsertLabelValues(db, table, labels)
	if err != nil {
		return fmt.Errorf("inserting labels for %s table: %w\n", table, err)
	}

	return nil
}

func PopulateMetadataTables(db *sql.DB) error {
	err := LoadMetadataTable(db, "../../resources/ages.csv", "ages")
	if err != nil {
		return fmt.Errorf("loading metadata from ages.csv: %w\n", err)
	}

	err = LoadMetadataTable(db, "../../resources/countries.csv", "countries")
	if err != nil {
		return fmt.Errorf("loading metadata from countries.csv: %w\n", err)
	}

	err = LoadMetadataTable(db, "../../resources/genders.csv", "genders")
	if err != nil {
		return fmt.Errorf("loading metadata from genders.csv: %w\n", err)
	}

	return nil
}

func InsertWeeklyDeathsData(records []eurostat.WeeklyDeathsRecord, db *sql.DB) error {
	tx, err := db.Begin()
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
	db *sql.DB,
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

	stmt, err := db.Prepare(WEEKLY_DEATHS_FOR_COUNTRY)
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

func parseLabelFile(path string) ([]MetadataLabel, error) {
	var records []MetadataLabel

	file, err := os.Open(path)
	if err != nil {
		return records, nil
	}

	r := csv.NewReader(file)

	rows, err := r.ReadAll()
	if err != nil {
		return records, nil
	}

	for _, row := range rows[1:] {
		value := row[0]
		label := row[1]
		order, err := strconv.Atoi(row[2])
		if err != nil {
			return records, err
		}
		records = append(records, MetadataLabel{
			Value: value,
			Label: label,
			Order: order,
		})
	}

	return records, nil
}

func InsertLabelValues(db *sql.DB, table string, labels []MetadataLabel) error {
	queryMap := map[string]string{
		"countries": COUNTRY_LABEL_INSERT_SQL,
		"ages":      AGE_LABEL_INSERT_SQL,
		"genders":   GENDER_LABEL_INSERT_SQL,
	}

	query, ok := queryMap[table]
	if !ok {
		return fmt.Errorf("no query found for %s table", table)
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("preparing %s query: %w\n", query, err)
	}

	for _, l := range labels {
		_, err = stmt.Exec(l.Value, l.Label, l.Order)
		if err != nil {
			return fmt.Errorf(
				"executing %s query for values (%s, %s, %d): %w\n",
				query,
				l.Value,
				l.Label,
				l.Order,
				err,
			)
		}
	}

	log.Printf("%d labels for table %s inserted.", len(labels), table)
	return nil
}

func GetLabels(db *sql.DB) ([]MetadataLabelFromDB, error) {
	var (
		value     string
		label     string
		order     int
		labelType string
	)
	var results []MetadataLabelFromDB

	stmt, err := db.Prepare(GET_LABLES_SQL)
	if err != nil {
		return results, fmt.Errorf("preparing %s query: %w\n", GET_LABLES_SQL, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return results, fmt.Errorf("executing %s query: %w\n", GET_LABLES_SQL, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&value, &label, &order, &labelType)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, MetadataLabelFromDB{
			Value:     value,
			Label:     label,
			Order:     order,
			LabelType: labelType,
		})
	}

	return results, nil
}