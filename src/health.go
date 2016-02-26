package main

import (
	"encoding/json"
	"net/http"
)

type Health struct {
	Up bool `json:"up"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	status := Health{true}

	json, err := json.Marshal(status)
	if err != nil {
		BroadcastError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
