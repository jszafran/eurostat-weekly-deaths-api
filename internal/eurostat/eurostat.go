package eurostat

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const DATA_URL = "https://ec.europa.eu/eurostat/estat-navtree-portlet-prod/BulkDownloadListing?file=data/demo_r_mwk_05.tsv.gz"

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
	Deaths  sql.NullInt64
	Age     string
	Gender  string
	Country string
}

// ReadData makes a HTTP requests to fetch gzipped TSV data from Eurostat website.
// Returns the TSV data as a string.
func ReadData() (string, error) {
	var data string

	log.Println("Fetching data from Eurostat.")
	resp, err := http.Get(DATA_URL)
	if err != nil {
		return data, fmt.Errorf("error when calling GET %s: %w\n", DATA_URL, err)
	}

	gzipBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return data, fmt.Errorf("opening gzip reader: %w\n", err)
	}

	tsvData, err := ioutil.ReadAll(gzipBody)
	if err != nil {
		return data, fmt.Errorf("reading gzip body: %w\n", err)
	}

	log.Println("Data fetched successfully.")
	return string(tsvData), nil
}

// WeekOfYearHeaderPositionMap calculates a mapping between header position index and week of year
// it represents.
func WeekOfYearHeaderPositionMap(header string) (map[int]WeekOfYear, error) {
	m := make(map[int]WeekOfYear)
	for i, v := range strings.Split(header, "\t")[1:] {
		woy, err := ParseWeekOfYear(v)
		if err != nil {
			return m, fmt.Errorf("parsing week of year for %s: %w\n", v, err)
		}
		m[i+1] = woy
	}
	return m, nil
}

// ParseMetadata parses a metadata information from a TSV line.
func ParseMetadata(line string) (Metadata, error) {
	var metadata Metadata

	meta := strings.Split(line, "\t")[0]
	parts := strings.Split(meta, ",")

	if len(parts) != 4 {
		return metadata, fmt.Errorf("parsing metadata: bad line metadata values %+v", parts)
	}
	return Metadata{
		Age:     parts[0],
		Gender:  parts[1],
		Country: parts[3],
	}, nil
}

// ParseDeathsValue parses information about reported amount of deaths.
// If no value was reported (or couldn't successfully parse the information),
// null value is returned (sql.NullInt64).
func ParseDeathsValue(v string) (sql.NullInt64, error) {
	var res sql.NullInt64
	v = strings.Replace(v, "p", "", -1)
	v = strings.Replace(v, ":", "", -1)
	v = strings.TrimSpace(v)

	i, err := strconv.Atoi(v)
	if err != nil {
		if v != "" {
			return res, fmt.Errorf("unparsable value %s: %w\n", v, err)
		}
		return sql.NullInt64{Valid: false}, nil
	}

	return sql.NullInt64{Int64: int64(i), Valid: true}, nil
}

// ParseWeekOfYear parses a week of year (WeekOfYear) information from given string.
func ParseWeekOfYear(s string) (WeekOfYear, error) {
	var woy WeekOfYear
	parts := strings.Split(strings.TrimSpace(s), "W")
	if len(parts) != 2 {
		return woy, fmt.Errorf("bad week of year value: %s", s)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return woy, fmt.Errorf("extracting year value from %s: %w\n", parts[0], err)
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return woy, fmt.Errorf("extracting week value from %s: %w\n", parts[1], err)
	}

	return WeekOfYear{
		Year: year,
		Week: week,
	}, nil
}

// ParseLine parses a full line ot TSV record. Returns a slice of WeeklyDeathsRecord
// or error (if unparsable data occurred).
func ParseLine(line string, woyPosMap map[int]WeekOfYear) ([]WeeklyDeathsRecord, error) {
	var wdr []WeeklyDeathsRecord

	metadata, err := ParseMetadata(line)
	if err != nil {
		return wdr, fmt.Errorf("extracting metadata from '%s': %w\n", line, err)
	}

	data := strings.Split(line, "\t")
	deaths := data[1:]

	for i, v := range deaths {
		dv, err := ParseDeathsValue(v)
		if err != nil {
			log.Fatalf("parsing deaths value %s: %s", v, err)
		}
		woy := woyPosMap[i+1]
		record := WeeklyDeathsRecord{
			Week:    woy.Week,
			Year:    woy.Year,
			Deaths:  dv,
			Age:     metadata.Age,
			Gender:  metadata.Gender,
			Country: metadata.Country,
		}
		wdr = append(wdr, record)
	}

	return wdr, nil
}

func ParseData(data string) ([]WeeklyDeathsRecord, error) {
	var records []WeeklyDeathsRecord

	split := strings.Split(data, "\n")
	header := split[0]
	rows := split[1:]

	woyPosMap, err := WeekOfYearHeaderPositionMap(header)
	if err != nil {
		return records, fmt.Errorf("creating week of year header position map: %w\n", err)
	}

	for i, line := range rows {
		parsedRecords, err := ParseLine(line, woyPosMap)
		if err != nil {
			return records, fmt.Errorf("parsing line no %d: %w\n", i, err)
		}

		for _, record := range parsedRecords {
			records = append(records, record)
		}
	}

	return records, nil
}