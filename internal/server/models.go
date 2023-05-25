package server

import "weekly_deaths/internal/eurostat"

type WeeklyDeathsResponse struct {
	Gender       string                    `json:"gender"`
	Age          string                    `json:"age"`
	Country      string                    `json:"country"`
	WeeklyDeaths []eurostat.WeekYearDeaths `json:"weekly_deaths"`
}
