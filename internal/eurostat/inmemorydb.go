package eurostat

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryDB struct {
	dataMu sync.RWMutex
	data   map[string][]WeeklyDeaths

	dataTimestampMu sync.RWMutex
	dataTimestamp   time.Time
}

func DBFromSnapshot(snapshot DataSnapshot) *InMemoryDB {
	return &InMemoryDB{
		data:          snapshot.Data,
		dataTimestamp: snapshot.Timestamp,
	}
}

func (db *InMemoryDB) GetWeeklyDeaths(
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
		key, err := makeKey(country, gender, age, year)
		if err != nil {
			return res, fmt.Errorf("fetching data from provider: %w", err)
		}

		db.dataMu.RLock()
		for _, r := range db.data[key] {
			res = append(res, WeekYearDeaths{Week: r.Week, Deaths: r.Deaths, Year: uint16(year)})
		}
		defer db.dataMu.RUnlock()
	}

	return res, nil
}

func (db *InMemoryDB) LoadSnapshot(snapshot DataSnapshot) {
	db.dataMu.Lock()
	db.dataTimestampMu.Lock()

	db.data = snapshot.Data
	db.dataTimestamp = snapshot.Timestamp

	defer db.dataMu.Unlock()
	defer db.dataTimestampMu.Unlock()
}

func (db *InMemoryDB) Timestamp() time.Time {
	db.dataTimestampMu.RLock()
	ts := db.dataTimestamp
	defer db.dataTimestampMu.RUnlock()

	return ts
}
