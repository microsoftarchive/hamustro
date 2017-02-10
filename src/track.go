package main

import (
	"bytes"
	"fmt"
	"github.com/bfaludi/remoteip"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/dialects"
	"github.com/wunderlist/hamustro/src/payload"
	"io/ioutil"
	"log"
	"mime"
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

	// Checks that the client want to send signature or not
	definedSignature := r.Header.Get("X-Hamustro-Time") != "" || r.Header.Get("X-Hamustro-Signature") != ""

	// If the client did not send time, we ignore
	if (signatureRequired || definedSignature) && r.Header.Get("X-Hamustro-Time") == "" {
		BroadcastError(w, "X-Hamustro-Time header is missing", http.StatusMethodNotAllowed)
		return
	}

	// If the client did not send signature of the message, we ignore
	if (signatureRequired || definedSignature) && r.Header.Get("X-Hamustro-Signature") == "" {
		BroadcastError(w, "X-Hamustro-Signature header is missing", http.StatusMethodNotAllowed)
		return
	}

	// Read the requests body into a variable.
	body, _ := ioutil.ReadAll(r.Body)

	// Calculate the request's signature
	if (signatureRequired || definedSignature) && r.Header.Get("X-Hamustro-Signature") != GetSignature(body, r.Header.Get("X-Hamustro-Time")) {
		BroadcastError(w, "X-Hamustro-Signature header is invalid", http.StatusMethodNotAllowed)
		return
	}

	collection := &payload.Collection{}
	contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch contentType {
	case "application/json":
		if err := jsonpb.Unmarshal(bytes.NewBuffer(body), collection); err != nil {
			BroadcastError(w, fmt.Sprintf("Unmarshaling json collection is failed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		if !collection.IsValid() {
			BroadcastError(w, fmt.Sprintf("Unmarshaled json collection is failed: required field not set"), http.StatusBadRequest)
			return
		}
	case "application/protobuf":
		if err := proto.Unmarshal(body, collection); err != nil {
			BroadcastError(w, fmt.Sprintf("Unmarshaling protobuf is failed: %s", err.Error()), http.StatusBadRequest)
			return
		}
	default:
		BroadcastError(w, "Unsupported or missing Content-Type", http.StatusBadRequest)
		return
	}

	// Checks the session information
	if GetSession(collection) != collection.GetSession() {
		BroadcastError(w, "Collection's session attribute is invalid", http.StatusBadRequest)
		return
	}

	// Stop if no payload information was received
	if !collection.HasPayloads() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Creates a Job and put into the JobQueue for processing.
	for _, payload := range collection.GetPayloads() {
		event := dialects.NewEvent(collection, payload)
		if event.IP == "" {
			if IP := remoteip.GetIPv4Address(r); IP != "" {
				event.SetIPAddress(IP)
			}
		}
		if config.IsMaskedIP() {
			event.TruncateIPv4LastOctet()
		}
		action := EventAction{event, 1}
		jobQueue <- &action
	}

	// Returns with 200.
	w.WriteHeader(http.StatusOK)
}
