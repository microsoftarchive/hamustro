package main

import (
	"encoding/json"
	"net/http"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"io"
)

type Flush struct {
	Flushed bool `json:"flushed"`
}

// Returns the request's signature
func GetMaintanceKey() string {
	maintanceKeyHash := sha256.New()
	io.WriteString(maintanceKeyHash, config.MaintanceKey)
	return hex.EncodeToString(maintanceKeyHash.Sum(nil))
}

func FlushHandler(w http.ResponseWriter, r *http.Request) {

	// Do not accept new events while the server is shutting down.
	if isTerminating {
		BroadcastError(w, "Server is currenly shutting down", http.StatusServiceUnavailable)
		return
	}

	// Ignore not GET messages.
	if r.Method != "GET" {
		BroadcastError(w, "Sending method is not GET", http.StatusMethodNotAllowed)
		return
	}

	// If the client did not send key of the message, we ignore
	if r.URL.Query().Get("maintance_key") == "" {
		BroadcastError(w, GetMaintanceKey(), http.StatusMethodNotAllowed)
		BroadcastError(w, "Maintance key is missing", http.StatusMethodNotAllowed)
		return
	}

	// Compare keys
	if r.URL.Query().Get("maintance_key") != GetMaintanceKey() {
		BroadcastError(w, "Maintance key is invalid", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Flushing workers")
	dispatcher.Flush()

	status := Flush{true}

	json, err := json.Marshal(status)
	if err != nil {
		BroadcastError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
