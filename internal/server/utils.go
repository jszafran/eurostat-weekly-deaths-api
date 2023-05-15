package server

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(statusCode int, w http.ResponseWriter, data any) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSONError(statusCode int, w http.ResponseWriter, errorMessage string) error {
	err := WriteJSON(statusCode, w, map[string]string{"error": errorMessage})
	if err != nil {
		return err
	}
	return nil
}