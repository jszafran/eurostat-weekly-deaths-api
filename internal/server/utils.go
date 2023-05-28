package server

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(statusCode int, w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
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
