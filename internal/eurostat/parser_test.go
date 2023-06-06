package eurostat

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestParseData(t *testing.T) {
	type TestRecord struct {
		key   string
		value []WeeklyDeaths
	}

	records := []TestRecord{
		{
			key: "AD|2021|TOTAL|F",
			value: []WeeklyDeaths{
				{Week: 1, Deaths: 1},
				{Week: 2, Deaths: 0},
				{Week: 3, Deaths: 0},
			},
		},
		{
			key: "PL|2021|TOTAL|T",
			value: []WeeklyDeaths{
				{Week: 1, Deaths: 0},
				{Week: 2, Deaths: 123},
				{Week: 3, Deaths: 212},
			},
		},
		{
			key: "GB|2021|TOTAL|M",
			value: []WeeklyDeaths{
				{Week: 1, Deaths: 0},
				{Week: 2, Deaths: 13},
				{Week: 3, Deaths: 25},
			},
		},
	}
	f, err := os.Open("testdata/eurostat_mockdata.tsv")
	if err != nil {
		log.Fatalf("Error when opening fixture data: %s\n", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		if err != nil {
			log.Fatalf("Error when reading fixture data: %s\n", err)
		}
	}

	parsedData, err := ParseData(string(data))
	if err != nil {
		log.Fatalf("Expected error to be nil but got %s\n", err)
	}

	if len(parsedData) == 0 {
		t.Fatal("Expected to get parsed records but received empty slice.")
	}

	for _, r := range records {
		got := parsedData[r.key]
		want := r.value
		if !reflect.DeepEqual(got, want) {
			log.Fatalf("Key %s: expected %+v but got %+v", r.key, want, got)
		}
	}
}
