package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

// Returns the request's signature
func GetMaintenanceKey() string {
	maintenanceKeyHash := sha256.New()
	io.WriteString(maintenanceKeyHash, config.MaintenanceKey)
	return hex.EncodeToString(maintenanceKeyHash.Sum(nil))
}

func FlushHandler(w http.ResponseWriter, r *http.Request) {
	// Do not accept flush request if
	if config.MaintenanceKey == "" {
		BroadcastError(w, "Please define maintanance key to access this feature", http.StatusServiceUnavailable)
		return
	}

	// Do not accept flush request while the server is shutting down.
	if isTerminating {
		BroadcastError(w, "Server is currenly shutting down", http.StatusServiceUnavailable)
		return
	}

	// Ignore not POST messages.
	if r.Method != "POST" {
		BroadcastError(w, "Sending method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// If the client did not send key of the message, we ignore
	if r.Header.Get("X-Hamustro-Maintenance-Key") == "" {
		BroadcastError(w, "Maintenance key is missing", http.StatusMethodNotAllowed)
		return
	}

	// Compare keys
	if r.Header.Get("X-Hamustro-Maintenance-Key") != GetMaintenanceKey() {
		BroadcastError(w, "Maintenance key is invalid", http.StatusMethodNotAllowed)
		return
	}

	log.Print("Flushing workers")
	dispatcher.Flush(&FlushOptions{Automatic: false})

	w.WriteHeader(http.StatusOK)
}
