package web

import (
	"time"

	"weekly_deaths/eurostat"
)

// WeeklyDeathsResponse represents a structure returned by
// /api/weekly_deaths endpoint.
type WeeklyDeathsResponse struct {
	Gender       string                    `json:"gender"`
	Age          string                    `json:"age"`
	Country      string                    `json:"country"`
	WeeklyDeaths []eurostat.WeekYearDeaths `json:"weekly_deaths"`
}

// MetadataLabel is a representation of label data
// that is returned by /api/labels endpoint.
type MetadataLabel struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Order int    `json:"order"`
	Type  string `json:"type"`
}

// InfoResponse is a representation of metadata info
// (hash of commit that application was built from and
// timestamp of downloading Eurostat data) returned by
// /api/info endpoint.
type InfoResponse struct {
	CommitHash       string    `json:"commit_hash"`
	DataDownloadedAt time.Time `json:"data_downloaded_at_utc_time"`
}
