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
