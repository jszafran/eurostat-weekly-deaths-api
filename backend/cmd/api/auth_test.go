package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const testUsername = "foo"
const testPassword = "bar"

var app = application{
	db: nil,
	auth: struct {
		username string
		password string
	}{
		username: testUsername,
		password: testPassword,
	},
}

var handler = app.basicAuth(TestingHandler)

func TestingHandler(w http.ResponseWriter, r *http.Request) {
	_ = writeJSON(http.StatusOK, w, map[string]string{"message": "hello"})
}

func TestEndpointProtectedWithBasicAuthWithFailedCredentials(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	req.SetBasicAuth(testUsername+"_", testPassword+"_")
	rr := httptest.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	wantStatus := http.StatusUnauthorized
	if sc := rr.Result().StatusCode; sc != wantStatus {
		t.Fatalf("Expected status %d but got %d\n", wantStatus, sc)
	}

	wantBody := []byte("Unauthorized\n")
	body, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("unexpected error when reading body: %s", body)
	}

	if !reflect.DeepEqual(wantBody, body) {
		t.Fatalf("expected %s but got %s", wantBody, body)
	}
}

func TestEndpointProtectedWithBasicAuthWithCorrectCredentials(t *testing.T) {
	var response struct {
		Message string `json:"message"`
	}

	req, err := http.NewRequest("GET", "", nil)
	req.SetBasicAuth(testUsername, testPassword)
	rr := httptest.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)
	wantStatus := http.StatusOK
	if sc := rr.Result().StatusCode; sc != wantStatus {
		t.Fatalf("Expected status %d but got %d\n", wantStatus, sc)
	}

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	wantMsg := "hello"
	if response.Message != wantMsg {
		t.Fatalf("Expected message %s but got %s\n", wantMsg, response.Message)
	}
}
