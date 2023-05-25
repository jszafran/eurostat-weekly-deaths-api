package eurostat

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/DmitriyVTitov/size"
)

var EurostatDataProvider *DataProvider

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

// DataProvider provides weekly deaths data fetched
// from Eurostat.
type DataProvider struct {
	data map[string][]WeeklyDeaths
}

// MakeKey creates a string key used for storing the data in
// application's memory (concatenation of country, gender, age and year).
func MakeKey(country string, gender string, age string, year int) (string, error) {
	yearStr := strconv.Itoa(year)
	if len(country) == 0 || len(gender) == 0 || len(age) == 0 || len(yearStr) == 0 {
		return "", fmt.Errorf("key cannot consist of empty string")
	}
	return fmt.Sprintf("%s|%d|%s|%s", country, year, age, gender), nil
}

// NewDataProvider returns a data provider object with
// data downloaded from provided data source.
func NewDataProvider(dataSource DataSource) (*DataProvider, error) {
	var dp DataProvider

	fmt.Println("Instantianting data provider.")
	rawData, err := dataSource.FetchData()
	if err != nil {
		return &dp, fmt.Errorf("fetching data from source: %w\n", err)
	}

	parsedData, err := ParseData(rawData)
	if err != nil {
		return &dp, fmt.Errorf("failed to parse raw data: %w\n", err)
	}

	for _, slice := range parsedData {
		sort.Slice(slice, func(i, j int) bool { return slice[i].Week < slice[j].Week })
	}

	log.Printf("Data provider instantiated (%d keys loaded to memory).\n", len(parsedData))
	log.Printf("Size of provider data: %d bytes", size.Of(parsedData))
	return &DataProvider{parsedData}, nil
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

// GetWeeklyDeaths returns weekly deaths data for given country,
// age, gender and year range.
func (dp *DataProvider) GetWeeklyDeaths(
	country string,
	age string,
	gender string,
	yearFrom int,
	yearTo int,
) ([]WeekYearDeaths, error) {
	res := make([]WeekYearDeaths, 0)

	years := makeRange(yearFrom, yearTo)
	if len(years) == 0 {
		return res, nil
	}

	for _, year := range years {
		key, err := MakeKey(country, gender, age, year)
		if err != nil {
			return res, fmt.Errorf("fetching data from provider: %w\n", err)
		}

		for _, r := range dp.data[key] {
			res = append(res, WeekYearDeaths{Week: r.Week, Deaths: r.Deaths, Year: uint16(year)})
		}
	}

	return res, nil
}
