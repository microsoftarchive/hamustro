package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Correct and incorrect header generation functions
type FlushHeaderFunction func() map[string]string

func GetMissingFlushHeader() map[string]string {
	return map[string]string{}
}
func GetFlushHeaderWithInvalidMaintenanceKey() map[string]string {
	return map[string]string{"X-Hamustro-Maintenance-Key": "fdsa43211"}
}
func GetValidFlushHeader() map[string]string {
	return map[string]string{"X-Hamustro-Maintenance-Key": GetMaintenanceKey()}
}

// Correct and incorrect config generation functions
type FlushConfigFunction func() *Config

func GetEmptyConfig() *Config {
	return &Config{}
}
func GetConfigWithMaintenanceKey() *Config {
	return &Config{MaintenanceKey: "maintancekey"}
}

// Input cases for the FlushHandler
type FlushHeaderTestCase struct {
	Method        string
	GetHeader     FlushHeaderFunction
	IsTerminating bool
	ExpectedCode  int
	GetConfig     FlushConfigFunction
}

// Executes the test cases for the given inputs
func RunTestsOnFlushHeader(t *testing.T, cases []*FlushHeaderTestCase) {
	for _, c := range cases {
		config = c.GetConfig()
		isTerminating = c.IsTerminating

		// Creates a new request
		req, _ := http.NewRequest(c.Method, "/api/flush", nil)

		// Set up the headers based on the predefined function
		for key, value := range c.GetHeader() {
			req.Header.Set(key, value)
		}
		resp := httptest.NewRecorder()

		FlushHandler(resp, req) // Calls the API

		// Check the status code
		if resp.Code != c.ExpectedCode {
			t.Errorf("Non-expected status code %d with the following body `%s`, it should be %d", resp.Code, resp.Body, c.ExpectedCode)
		}
	}
}

// Tests the API
func TestFlushHeader(t *testing.T) {
	t.Log("Test flush header")
	config = &Config{}                                           // Creates a config
	storageClient = &BufferedStorageClient{}                     // Define the Buffered Storage as a storage
	jobQueue = make(chan Job, 10)                                // Creates a jobQueue
	log.SetOutput(ioutil.Discard)                                // Disable the logger
	T, response, catched = t, nil, false                         // Set properties for the BufferedStorageClient
	dispatcher = NewDispatcher(1, &WorkerOptions{BufferSize: 5}) // Creates a dispatcher
	dispatcher.Run()

	t.Log("Test flush headers with different setups")
	RunTestsOnFlushHeader(t,
		[]*FlushHeaderTestCase{
			{"POST", GetValidFlushHeader, false, http.StatusServiceUnavailable, GetEmptyConfig},
			{"POST", GetValidFlushHeader, true, http.StatusServiceUnavailable, GetConfigWithMaintenanceKey},
			{"GET", GetValidFlushHeader, false, http.StatusMethodNotAllowed, GetConfigWithMaintenanceKey},
			{"POST", GetMissingFlushHeader, false, http.StatusMethodNotAllowed, GetConfigWithMaintenanceKey},
			{"POST", GetFlushHeaderWithInvalidMaintenanceKey, false, http.StatusMethodNotAllowed, GetConfigWithMaintenanceKey},
			{"POST", GetValidFlushHeader, false, http.StatusOK, GetConfigWithMaintenanceKey},
		})

}
