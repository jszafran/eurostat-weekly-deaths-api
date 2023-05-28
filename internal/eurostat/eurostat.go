package eurostat

import (
	"compress/gzip"
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
	Deaths  int
	Age     string
	Gender  string
	Country string
}

type DataSource interface {
	FetchData() (string, error)
}

type LiveEurostatDataSource struct{}

func (s LiveEurostatDataSource) FetchData() (string, error) {
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
// 0 is returned.
func ParseDeathsValue(v string) (int, error) {
	var res int
	v = strings.Replace(v, "p", "", -1)
	v = strings.Replace(v, ":", "", -1)
	v = strings.TrimSpace(v)

	i, err := strconv.Atoi(v)
	if err != nil {
		if v != "" {
			return res, fmt.Errorf("unparsable value %s: %w\n", v, err)
		}
		return 0, nil
	}

	return i, nil
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
func ParseLine(line string, woyPosMap map[int]WeekOfYear, results map[string][]WeeklyDeaths) error {
	metadata, err := ParseMetadata(line)
	if err != nil {
		return fmt.Errorf("extracting metadata from '%s': %w\n", line, err)
	}

	data := strings.Split(line, "\t")
	deaths := data[1:]

	for i, v := range deaths {
		dv, err := ParseDeathsValue(v)
		if err != nil {
			log.Fatalf("parsing deaths value %s: %s", v, err)
		}
		woy := woyPosMap[i+1]
		key, err := MakeKey(metadata.Country, metadata.Gender, metadata.Age, woy.Year)
		if err != nil {
			return fmt.Errorf("failed to create key for %+v metadata and %+v week of year\n", metadata, woy)
		}

		// Year, according to ISO definitions, contains
		// 52 or 53 full weeks. Eurostat dataset contains
		// column with week=99, hence below condition
		// filtering them out.
		if woy.Week >= 54 {
			continue
		}

		results[key] = append(results[key], WeeklyDeaths{Week: uint8(woy.Week), Deaths: uint32(dv)})

	}

	return nil
}

func ParseData(data string) (map[string][]WeeklyDeaths, error) {
	results := make(map[string][]WeeklyDeaths)

	split := strings.Split(data, "\n")
	header := split[0]
	rows := split[1:]

	woyPosMap, err := WeekOfYearHeaderPositionMap(header)
	if err != nil {
		return nil, fmt.Errorf("creating week of year header position map: %w\n", err)
	}

	for i, line := range rows {
		if line == "" {
			continue
		}
		err := ParseLine(line, woyPosMap, results)
		if err != nil {
			return results, fmt.Errorf("parsing line no %d: %w\n", i, err)
		}
	}

	return results, nil
}

// UniqueValues return unique values for age, gender, country
// derived from parsed Eurostat data.
func uniqueValues(data []WeeklyDeathsRecord) map[string][]string {
	var v interface{}
	result := make(map[string][]string)

	ages := make(map[string]interface{})
	genders := make(map[string]interface{})
	countries := make(map[string]interface{})

	for _, rec := range data {
		ages[rec.Age] = v
		genders[rec.Gender] = v
		countries[rec.Country] = v
	}

	for k := range ages {
		result["age"] = append(result["age"], k)
	}

	for k := range genders {
		result["gender"] = append(result["gender"], k)
	}

	for k := range countries {
		result["country"] = append(result["country"], k)
	}

	return result
}

func contains(container []string, value string) bool {
	for _, el := range container {
		if value == el {
			return true
		}
	}

	return false
}

func missingLabels(data []WeeklyDeathsRecord) map[string][]string {
	res := make(map[string][]string)
	uv := uniqueValues(data)

	fixedAgeLabels := make([]string, 0)
	for _, age := range ageLabels {
		fixedAgeLabels = append(fixedAgeLabels, age.Value)
	}

	for _, a := range uv["age"] {
		if !contains(fixedAgeLabels, a) {
			res["age"] = append(res["age"], a)
		}
	}

	fixedCountryLabels := make([]string, 0)
	for _, country := range countryLabels {
		fixedCountryLabels = append(fixedCountryLabels, country.Value)
	}

	for _, c := range uv["country"] {
		if !contains(fixedCountryLabels, c) {
			res["country"] = append(res["country"], c)
		}
	}

	fixedGenderLabels := make([]string, 0)
	for _, gender := range genderLabels {
		fixedGenderLabels = append(fixedGenderLabels, gender.Value)
	}

	for _, g := range uv["gender"] {
		if !contains(fixedGenderLabels, g) {
			res["gender"] = append(res["gender"], g)
		}
	}

	return res
}

// ValidateLabels compares the fixed country, age, gender labels against
// data fetched from Eurostat. Reports error in case of any discrepancies.
func ValidateLabels(data []WeeklyDeathsRecord) error {
	missingDataFound := false
	ml := missingLabels(data)
	age, exists := ml["age"]
	if exists {
		log.Printf("Missing age labels found: %+v\n", age)
		missingDataFound = true
	}

	gender, exists := ml["gender"]
	if exists {
		log.Printf("Missing gender labels found: %+v\n", gender)
		missingDataFound = true
	}

	country, exists := ml["country"]
	if exists {
		log.Printf("Missing country labels found: %+v\n", country)
		missingDataFound = true
	}

	if missingDataFound {
		return fmt.Errorf("found missing labels: %+v\n", ml)
	}

	return nil
}
