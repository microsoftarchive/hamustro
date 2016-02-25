package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/dialects"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func GetTestProtobufCollectionBody(userId uint32, disturbSession bool, numberOfPayloads int) ([]byte, []*Job) {
	collection := GetTestPayloadCollection(userId, numberOfPayloads)
	if disturbSession {
		collection.Session = proto.String("not-valid-session")
	}
	data, _ := proto.Marshal(collection)
	var jobs []*Job
	for _, payload := range collection.GetPayloads() {
		jobs = append(jobs, &Job{dialects.NewEvent(collection, payload), 1})
	}
	return data, jobs
}

func TestTrackHandler(t *testing.T) {
	t.Log("Creating new workers")
	storageClient = &SimpleStorageClient{}                          // Define the Simple Storage as a storage
	jobQueue = make(chan *Job, 10)                                  // Creates a jobQueue
	log.SetOutput(ioutil.Discard)                                   // Disable the logger
	T, response, catched = t, nil, false                            // Set properties for the SimpleStorageClient
	dispatcher := NewDispatcher(2, &WorkerOptions{RetryAttempt: 5}) // Creates a dispatcher
	dispatcher.Run()                                                // Starts the dispatcher
	config = &Config{SharedSecret: "ultrasafesecret"}               // Creates a config

	if exp := 2; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	rTime := "1454514088"
	var (
		notValidBody      = []byte("orange")
		notValidSignature = GetSignature(notValidBody, rTime)
	)
	var (
		almostValidBody, _   = GetTestProtobufCollectionBody(2423, true, 1)
		almostValidSignature = GetSignature(almostValidBody, rTime)
	)
	var (
		validBody, validJobs = GetTestProtobufCollectionBody(633289, false, 1)
		validSignature       = GetSignature(validBody, rTime)
	)
	var (
		mValidBody, mValidJobs = GetTestProtobufCollectionBody(53464, false, 2)
		mValidSignature        = GetSignature(mValidBody, rTime)
	)
	cases := []struct {
		Method        string
		Body          []byte
		Header        map[string]string
		IsTerminating bool
		ExpectedCode  int
		ExpectedJobs  []*Job
	}{
		{"GET", notValidBody, map[string]string{}, true, http.StatusServiceUnavailable, nil},                                                                          // 1. Service is shutting down
		{"GET", notValidBody, map[string]string{}, false, http.StatusMethodNotAllowed, nil},                                                                           // 2. GET is not supported
		{"POST", notValidBody, map[string]string{}, false, http.StatusMethodNotAllowed, nil},                                                                          // 3. Missing headers
		{"POST", notValidBody, map[string]string{"X-Hamustro-Signature": notValidSignature}, false, http.StatusMethodNotAllowed, nil},                                 // 4. Missing X-Hamustro-rTime
		{"POST", notValidBody, map[string]string{"X-Hamustro-Time": rTime}, false, http.StatusMethodNotAllowed, nil},                                                  // 5. Missing X-Hamustro-Signature
		{"POST", notValidBody, map[string]string{"X-Hamustro-Signature": notValidSignature + "x", "X-Hamustro-Time": rTime}, false, http.StatusMethodNotAllowed, nil}, // 6. X-Hamustro-Signature is invalid
		{"POST", notValidBody, map[string]string{"X-Hamustro-Signature": notValidSignature, "X-Hamustro-Time": rTime}, false, http.StatusBadRequest, nil},             // 7. Content is not valid protobuf
		{"POST", almostValidBody, map[string]string{"X-Hamustro-Signature": almostValidSignature, "X-Hamustro-Time": rTime}, false, http.StatusBadRequest, nil},       // 8. Session is not valid
		{"POST", validBody, map[string]string{"X-Hamustro-Signature": validSignature, "X-Hamustro-Time": rTime}, false, http.StatusOK, validJobs},                     // 9. This message is valid
		{"POST", mValidBody, map[string]string{"X-Hamustro-Signature": mValidSignature, "X-Hamustro-Time": rTime}, false, http.StatusOK, mValidJobs},                  // 10. This message is valid and has 3 subjob
	}

	for i, c := range cases {
		isTerminating = c.IsTerminating
		for _, isVerbose := range []bool{true, false} {
			verbose = isVerbose         // Sets the verbose mode
			exp = map[string]struct{}{} // Resets the expectations dict

			t.Logf("Checking the %d test case with %s mode", i+1, verbose)

			req, _ := http.NewRequest(c.Method, "/api/v1/track", bytes.NewBuffer(c.Body))
			for key, value := range c.Header {
				req.Header.Set(key, value)
			}
			resp := httptest.NewRecorder()

			if c.ExpectedJobs != nil {
				for _, job := range c.ExpectedJobs {
					SetJobExpectation([]*Job{job}, false, false)
				}
			}
			TrackHandler(resp, req)
			if c.ExpectedJobs != nil {
				time.Sleep(150 * time.Millisecond)
				ValidateSending()
			}

			if resp.Code != c.ExpectedCode {
				t.Errorf("Non-expected status code %d with the following body `%s`, it should be %d", resp.Code, resp.Body, c.ExpectedCode)
			}
		}
	}

	dispatcher.Stop()
}
