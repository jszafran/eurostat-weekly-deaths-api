package queries

const CREATE_WEEKLY_DEATHS_SQL = `
	CREATE TABLE weekly_deaths (
		week INTEGER NOT NULL,
		year INTEGER NOT NULL,
		deaths INTEGER,
		age STRING,
		gender STRING,
		country STRING
	) 
`

const CREATE_COUNTRIES_SQL = `CREATE TABLE countries (value STRING, label STRING, "order" INTEGER)`

const CREATE_AGES_SQL = `CREATE TABLE ages (value STRING, label STRING, "order" INTEGER)`

const CREATE_GENDERS_SQL = `CREATE TABLE genders (value STRING, label STRING, "order" INTEGER)`

const WEEKLY_DEATHS_FOR_COUNTRY = `
	SELECT
		week,
		year,
		deaths,
		age,
		gender,
		country
	FROM weekly_deaths
	WHERE 1=1
	AND country = ?
	AND gender = ?
	AND age = ?
	AND year >= ?
	AND year <= ?
`

const COUNTRY_LABEL_INSERT_SQL = `INSERT INTO countries (value, label, "order") VALUES (?, ?, ?)`

const GENDER_LABEL_INSERT_SQL = `INSERT INTO genders (value, label, "order") VALUES (?, ?, ?)`

const AGE_LABEL_INSERT_SQL = `INSERT INTO ages (value, label, "order") VALUES (?, ?, ?)`
