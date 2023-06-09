package eurostat

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	maxIsoWeekNum = 53

	// metadata column should contain 4 elements
	// after splitting by coma
	metadataElementsLength = 4
	// week year value should contain 2 elements
	// after splitting by W character
	weekYearElementsLength = 2

	tabulator = string(rune(0x09))
)

// WeekOfYear represents a single week of year (ISO week).
type weekOfYear struct {
	Year int
	Week int
}

// ParseData parses Eurostat raw string data into key value data store,
// where key is a combination of country, age, gender and year values
// and value is a slice of WeeklyDeaths struct.
func ParseData(r io.ReadCloser) (map[string][]WeeklyDeaths, error) {
	log.Println("Starting parsing data.")
	results := make(map[string][]WeeklyDeaths)

	data, err := io.ReadAll(r)
	if err != nil {
		return results, err
	}
	defer r.Close()

	split := strings.Split(string(data), "\n")
	header := split[0]
	rows := split[1:]

	woyPosMap, err := weekOfYearHeaderPositionMap(header)
	if err != nil {
		return nil, fmt.Errorf("creating week of year header position map: %w", err)
	}

	for i, line := range rows {
		if line == "" {
			continue
		}
		err := parseLine(line, woyPosMap, results)
		if err != nil {
			return results, fmt.Errorf("parsing line no %d: %w", i, err)
		}
	}

	log.Println("Parsing finished.")

	for k := range results {
		sort.Slice(results[k], func(i, j int) bool {
			return results[k][i].Week < results[k][j].Week
		})
	}
	return results, nil
}

func parseWeekOfYear(s string) (weekOfYear, error) {
	var woy weekOfYear
	parts := strings.Split(strings.TrimSpace(s), "W")

	if len(parts) != weekYearElementsLength {
		return woy, fmt.Errorf("bad week of year value: %s", s)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return woy, fmt.Errorf("extracting year value from %s: %w", parts[0], err)
	}

	week, err := strconv.Atoi(parts[1])
	if err != nil {
		return woy, fmt.Errorf("extracting week value from %s: %w", parts[1], err)
	}

	return weekOfYear{
		Year: year,
		Week: week,
	}, nil
}

// ParseDeathsValue parses information about reported amount of deaths.
// If no value was reported (or couldn't successfully parse the information),
// 0 is returned.
func parseDeathsValue(v string) (int, error) {
	var res int
	v = strings.Replace(v, "p", "", -1)
	v = strings.Replace(v, ":", "", -1)
	v = strings.TrimSpace(v)

	i, err := strconv.Atoi(v)
	if err != nil {
		if v != "" {
			return res, fmt.Errorf("unparsable value %s: %w", v, err)
		}
		return 0, nil
	}

	return i, nil
}

func parseMetadata(line string) (Metadata, error) {
	var metadata Metadata

	meta := strings.Split(line, tabulator)[0]
	parts := strings.Split(meta, ",")

	if len(parts) != metadataElementsLength {
		return metadata, fmt.Errorf("parsing metadata: bad line metadata values %+v", parts)
	}
	return Metadata{
		Age:     parts[0],
		Gender:  parts[1],
		Country: parts[3],
	}, nil
}

func weekOfYearHeaderPositionMap(header string) (map[int]weekOfYear, error) {
	m := make(map[int]weekOfYear)
	for i, v := range strings.Split(header, tabulator)[1:] {
		woy, err := parseWeekOfYear(v)
		if err != nil {
			return m, fmt.Errorf("parsing week of year for %s: %w", v, err)
		}
		m[i+1] = woy
	}
	return m, nil
}

func parseLine(line string, woyPosMap map[int]weekOfYear, results map[string][]WeeklyDeaths) error {
	metadata, err := parseMetadata(line)
	if err != nil {
		return fmt.Errorf("extracting metadata from '%s': %w", line, err)
	}

	data := strings.Split(line, tabulator)
	deaths := data[1:]

	for i, v := range deaths {
		dv, err := parseDeathsValue(v)
		if err != nil {
			log.Fatalf("parsing deaths value %s: %s", v, err)
		}
		woy := woyPosMap[i+1]
		key, err := makeKey(metadata.Country, metadata.Gender, metadata.Age, woy.Year)
		if err != nil {
			return fmt.Errorf("failed to create key for %+v metadata and %+v week of year", metadata, woy)
		}

		// Year, according to ISO definitions, contains
		// 52 or 53 full weeks. Eurostat dataset contains
		// column with week=99, hence below condition
		// filtering them out.
		if woy.Week >= maxIsoWeekNum {
			continue
		}

		results[key] = append(results[key], WeeklyDeaths{Week: uint8(woy.Week), Deaths: uint32(dv)})
	}

	return nil
}
