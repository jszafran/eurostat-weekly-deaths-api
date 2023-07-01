package eurostat

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	eurostatDataUrl = "https://ec.europa.eu/eurostat/estat-navtree-portlet-prod/BulkDownloadListing?file=data/demo_r_mwk_05.tsv.gz"

	timestampLayout   = "20060102T150405"
	dataFileExtension = ".tsv.gz"
)

// WeeklyDeaths represents a number of deaths reported
// WeeklyDeaths represents a number of deaths reported
// for given week. Lack of information is represented
// with -1 value.
type WeeklyDeaths struct {
	Week   uint8  `json:"week"`
	Deaths uint32 `json:"deaths"`
}

type WeekYearDeaths struct {
	Week   uint8  `json:"week"`
	Year   uint16 `json:"year"`
	Deaths uint32 `json:"deaths"`
}

// WeekOfYear represents a single week of year (ISO week).
type WeekOfYear struct {
	Year int
	Week int
}

// Metadata contains information about age, gender and country of particular record.
type Metadata struct {
	Age     string
	Gender  string
	Country string
}

// WeeklyDeathsRecord represents a full record of weekly deaths data:
// - week
// - year
// - number of deaths (if reported)
// - age bucket
// - gender
// - country
type WeeklyDeathsRecord struct {
	Week    int
	Year    int
	Deaths  int
	Age     string
	Gender  string
	Country string
}

type DataSnapshot struct {
	Data      map[string][]WeeklyDeaths
	Timestamp time.Time
}

// TODO: implement S3 data provider

func textFromGZIP(r io.Reader) (string, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}

	text, err := io.ReadAll(gzipReader)
	if err != nil {
		return "", err
	}

	return string(text), nil
}

// makeKey creates a string key used for storing the data in
// application's memory (concatenation of country, gender, age and year).
func makeKey(country string, gender string, age string, year int) (string, error) {
	yearStr := strconv.Itoa(year)
	if len(country) == 0 || len(gender) == 0 || len(age) == 0 || len(yearStr) == 0 {
		return "", errors.New("key cannot consist of empty string")
	}
	return fmt.Sprintf("%s|%d|%s|%s", country, year, age, gender), nil
}

func makeRange(from int, to int) []int {
	rng := make([]int, 0)
	if from > to {
		return rng
	}

	if from == to {
		rng = append(rng, from)
		return rng
	}

	for i := from; i <= to; i++ {
		rng = append(rng, i)
	}

	return rng
}

func parseTimestamp(name string) (time.Time, error) {
	var ts time.Time
	name = strings.Replace(name, dataFileExtension, "", -1)
	ts, err := time.Parse(timestampLayout, name)
	if err != nil {
		return ts, err
	}
	return ts, nil
}

func timestampFromFileName(filePath string) (time.Time, error) {
	var ts time.Time

	filePath = path.Base(filePath)
	ts, err := parseTimestamp(filePath)
	if err != nil {
		return ts, err
	}

	return ts, nil
}

func DataSnapshotFromPath(path string) (DataSnapshot, error) {
	var ds DataSnapshot
	file, err := os.Open(path)
	if err != nil {
		return ds, err
	}

	r, err := gzip.NewReader(file)
	if err != nil {
		return ds, err
	}

	parsedData, err := ParseData(r)
	if err != nil {
		return ds, err
	}

	ts, err := timestampFromFileName(path)
	if err != nil {
		return ds, err
	}

	ds.Data = parsedData
	ds.Timestamp = ts
	return ds, nil
}

func persistSnapshot(b []byte, timestamp time.Time) {
	if os.Getenv("PERSIST_LIVE_SNAPSHOTS") != "true" {
		return
	}

	smg, err := NewSnapshotManager(os.Getenv("S3_BUCKET"))
	if err != nil {
		log.Printf("Failed to create snapshot manager: %s", err)
		return
	}

	err = smg.PersistSnapshot(bytes.NewReader(b), timestamp)
	if err != nil {
		log.Printf("Failed to persist snapshot to S3: %s", err)
		return
	}

	log.Println("Snapshot successfully persisted to S3!")
}

func DataSnapshotFromEurostat() (DataSnapshot, error) {
	var ds DataSnapshot

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", eurostatDataUrl, nil)
	if err != nil {
		return ds, err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return ds, err
	}
	defer resp.Body.Close()

	timestamp := time.Now().UTC()
	ds.Timestamp = timestamp

	sb, err := io.ReadAll(resp.Body)
	if err != nil {
		return ds, err
	}

	r, err := gzip.NewReader(bytes.NewReader(sb))
	if err != nil {
		return ds, err
	}

	data, err := ParseData(r)
	if err != nil {
		return ds, err
	}
	ds.Data = data

	persistSnapshot(sb, timestamp)
	return ds, nil
}
