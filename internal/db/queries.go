package db

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

const DROP_WEEKLY_DEATHS_SQL = `DROP TABLE IF EXISTS weekly_deaths`

const SELECT_WEEKLY_DEATHS_FOR_COUNTRY = `
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

const CREATE_INDEX_SQL = `
	CREATE INDEX idx_weekly_deaths 
	ON weekly_deaths (country, gender, age, year)
`

const DROP_INDEX_SQL = `DROP INDEX IF EXISTS idx_weekly_deaths`

const DELETE_INCORRECT_WEEKS_DATA_SQL = `DELETE FROM weekly_deaths WHERE week > 53`
