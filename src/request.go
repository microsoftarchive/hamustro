package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/dialects"
	"github.com/wunderlist/hamustro/src/payload"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// Prints the error messages.
func BroadcastError(w http.ResponseWriter, err string, code int) {
	log.Println(err)
	if verbose {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(code)
	if verbose {
		fmt.Fprintf(w, `{"error":%q}`, err)
	}
}

// Controller for `/api/v1/track`
func TrackHandler(w http.ResponseWriter, r *http.Request) {
	// Do not accept new events while the server is shutting down.
	if isTerminating {
		BroadcastError(w, "Server is currenly shutting down", http.StatusServiceUnavailable)
		return
	}

	// Ignore not POST messages.
	if r.Method != "POST" {
		BroadcastError(w, "Sending method is not POST", http.StatusMethodNotAllowed)
		return
	}

	// If the client did not send time, we ignore
	if r.Header.Get("X-Hamustro-Time") == "" {
		BroadcastError(w, "X-Hamustro-Time header is missing", http.StatusMethodNotAllowed)
		return
	}

	// If the client did not send signature of the message, we ignore
	if r.Header.Get("X-Hamustro-Signature") == "" {
		BroadcastError(w, "X-Hamustro-Signature header is missing", http.StatusMethodNotAllowed)
		return
	}

	// Read the requests body into a variable.
	body, _ := ioutil.ReadAll(r.Body)

	// Calculate the request's signature
	if r.Header.Get("X-Hamustro-Signature") != GetSignature(body, r.Header.Get("X-Hamustro-Time")) {
		BroadcastError(w, "X-Hamustro-Signature header is invalid", http.StatusMethodNotAllowed)
		return
	}

	// Read the body into protobuf decoding.
	collection := &payload.Collection{}
	if err := proto.Unmarshal(body, collection); err != nil {
		BroadcastError(w, fmt.Sprintf("Unmarshaling protobuf collection is failed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Checks the session information
	if GetSession(collection) != collection.GetSession() {
		BroadcastError(w, "Collection's session attribute is invalid", http.StatusBadRequest)
		return
	}

	// Creates a Job and put into the JobQueue for processing.
	for _, payload := range collection.Payloads {
		job := Job{dialects.NewEvent(collection, payload), 1}
		jobQueue <- &job
	}

	// Returns with 200.
	w.WriteHeader(http.StatusOK)
}
