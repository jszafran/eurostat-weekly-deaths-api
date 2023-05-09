package server

import "database/sql"

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
