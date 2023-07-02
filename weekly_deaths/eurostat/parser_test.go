package eurostat

import (
	"bytes"
	"compress/gzip"
	"log"
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
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte(`age,sex,unit,geo\time	2021W03	2021W02	2021W01
TOTAL,F,NR,AD	:	:	1
TOTAL,T,NR,PL	212	123	:
TOTAL,M,NR,GB	25 p	13	:`))

	w.Close()

	r, err := gzip.NewReader(&buf)
	if err != nil {
		log.Fatal(err)
	}

	parsedData, err := ParseData(r)
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
