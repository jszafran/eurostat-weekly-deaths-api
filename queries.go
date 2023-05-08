package main

const CREATE_WEEKLY_DEATHS_SQL = `
	CREATE TABLE weekly_deaths (
		week INTEGER NOT NULL,
		year INTEGER NOT NULL,
		deaths INTEGER,
		age STRING,
		sex STRING,
		country STRING
	) 
`

const CREATE_COUNTRIES_SQL = `CREATE TABLE countries (name string)`

const CREATE_AGES_SQL = `CREATE TABLE ages (name string)`

const CREATE_GENDERS_SQL = `CREATE TABLE genders (name string)`
