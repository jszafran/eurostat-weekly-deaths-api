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
