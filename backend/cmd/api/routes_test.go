package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"weekly_deaths/internal/eurostat"
)

const expectedContentType = "application/json"

func testTimestamp() time.Time {
	return time.Date(2021, 1, 12, 10, 23, 11, 0, time.UTC)
}

func testingDB() *eurostat.InMemoryDB {
	snapshot := eurostat.DataSnapshot{
		Data: map[string][]eurostat.WeeklyDeaths{
			"PL|2020|TOTAL|T": {
				{Week: 1, Deaths: 0},
				{Week: 2, Deaths: 0},
				{Week: 3, Deaths: 0},
				{Week: 4, Deaths: 1},
			},
			"PL|2021|TOTAL|T": {
				{Week: 1, Deaths: 5},
				{Week: 2, Deaths: 10},
				{Week: 3, Deaths: 15},
				{Week: 4, Deaths: 20},
			},
			"PL|2022|TOTAL|T": {
				{Week: 1, Deaths: 25},
				{Week: 2, Deaths: 30},
				{Week: 3, Deaths: 35},
				{Week: 4, Deaths: 40},
			},
			"GB|2012|TOTAL|F": {
				{Week: 1, Deaths: 100},
				{Week: 2, Deaths: 200},
				{Week: 3, Deaths: 300},
				{Week: 4, Deaths: 400},
			},
		},
		Timestamp: testTimestamp(),
	}
	return eurostat.DBFromSnapshot(snapshot)
}

func TestInfoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	commit := "6e874a04a4ebeb82128e2b2000c97649028218b6"
	os.Setenv("COMMIT", commit)

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.InfoHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	expectedTimestamp := "2021-01-12T10:23:11Z"
	expectedBody := fmt.Sprintf(
		"{\"commit_hash\":\"%s\",\"data_downloaded_at_utc_time\":\"%s\"}",
		commit,
		expectedTimestamp,
	)

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %d want %d", status, expectedStatus)
	}

	if body := strings.TrimSuffix(rr.Body.String(), "\n"); body != expectedBody {
		t.Errorf("handler returned unexpected body: got %s want %s", body, expectedBody)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
	}
}

func TestWeeklyDeathsHandlerWithIncorrectYearRange(t *testing.T) {
	var resp WeeklyDeathsResponse

	req, err := http.NewRequest("GET", "?country=PL&age=TOTAL&gender=T&year_from=2021&year_to=2020", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.WeeklyDeathsHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.WeeklyDeaths) > 0 {
		t.Errorf("handler returned unexpected body (non-empty slice): %+v\n", resp.WeeklyDeaths)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
	}
}

func TestWeeklyDeathsHandlerFetchingSingleYearData(t *testing.T) {
	var resp WeeklyDeathsResponse

	req, err := http.NewRequest("GET", "?country=PL&age=TOTAL&gender=T&year_from=2021&year_to=2021", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.WeeklyDeathsHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	want := []eurostat.WeekYearDeaths{
		{Week: 1, Year: 2021, Deaths: 5},
		{Week: 2, Year: 2021, Deaths: 10},
		{Week: 3, Year: 2021, Deaths: 15},
		{Week: 4, Year: 2021, Deaths: 20},
	}

	got := resp.WeeklyDeaths
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("handler returned unexpected body: want %+v but got %+v\n", want, got)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
	}
}

func TestWeeklyDeathsHandlerFetchingDataForRangeOfYears(t *testing.T) {
	var resp WeeklyDeathsResponse

	req, err := http.NewRequest("GET", "?country=PL&age=TOTAL&gender=T&year_from=2020&year_to=2022", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.WeeklyDeathsHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	want := []eurostat.WeekYearDeaths{
		{Week: 1, Year: 2020, Deaths: 0},
		{Week: 2, Year: 2020, Deaths: 0},
		{Week: 3, Year: 2020, Deaths: 0},
		{Week: 4, Year: 2020, Deaths: 1},
		{Week: 1, Year: 2021, Deaths: 5},
		{Week: 2, Year: 2021, Deaths: 10},
		{Week: 3, Year: 2021, Deaths: 15},
		{Week: 4, Year: 2021, Deaths: 20},
		{Week: 1, Year: 2022, Deaths: 25},
		{Week: 2, Year: 2022, Deaths: 30},
		{Week: 3, Year: 2022, Deaths: 35},
		{Week: 4, Year: 2022, Deaths: 40},
	}

	got := resp.WeeklyDeaths
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("handler returned unexpected body: want %+v but got %+v\n", want, got)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
	}
}

func TestWeeklyDeathsHandlerFetchingDataForNonexistingKey(t *testing.T) {
	var resp WeeklyDeathsResponse

	req, err := http.NewRequest("GET", "?country=DE&age=TOTAL&gender=T&year_from=2019&year_to=2022", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.WeeklyDeathsHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	got := resp.WeeklyDeaths
	if len(got) != 0 {
		log.Fatalf("handler returned unexpected body: wanted empty slice but got %+v\n", got)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
	}
}

func TestWeeklyDeathsHandlerMissingQueryParams(t *testing.T) {
	type ErroResponse struct {
		ErrorMessage string `json:"error"`
	}

	type TestCase struct {
		queryParams          string
		expectedErrorMessage string
	}

	var resp ErroResponse

	testCases := []TestCase{
		{queryParams: "?", expectedErrorMessage: "country url param required"},
		{queryParams: "?country=PL", expectedErrorMessage: "gender url param required"},
		{queryParams: "?country=PL&gender=T", expectedErrorMessage: "age url param required"},
		{queryParams: "?country=PL&gender=T&age=TOTAL", expectedErrorMessage: "year_from url param required"},
		{queryParams: "?country=PL&gender=T&age=TOTAL&year_from=2020", expectedErrorMessage: "year_to url param required"},
	}

	app := application{db: testingDB()}
	handler := http.HandlerFunc(app.WeeklyDeathsHandler)

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", tc.queryParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		json.NewDecoder(rr.Body).Decode(&resp)

		if resp.ErrorMessage != tc.expectedErrorMessage {
			t.Fatalf("expected %s but got %s", tc.expectedErrorMessage, resp.ErrorMessage)
		}

		if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
			t.Errorf("handler returned unexpected content-type: got %s want %s", contentType, expectedContentType)
		}
	}
}
