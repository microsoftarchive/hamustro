package main

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/dialects"
	"github.com/wunderlist/hamustro/src/payload"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Input collections for the test cases
type TrackBodyCollection struct {
	Collection               []byte
	DisturbedCollection      []byte
	IncompleteCollectionBody []byte
}

// Correct and incorrect header generation functions
type BodyFunction func(t *TrackBodyCollection) []byte

func GetCollectionBody(t *TrackBodyCollection) []byte {
	return t.Collection
}
func GetDisturbedCollectionBody(t *TrackBodyCollection) []byte {
	return t.DisturbedCollection
}
func GetIncompleteCollectionBody(t *TrackBodyCollection) []byte {
	return t.IncompleteCollectionBody
}

// Returns jobs from a collection (used for expected output)
func GetJobsFromCollection(collection *payload.Collection) []*EventAction {
	var jobs []*EventAction
	for _, payload := range collection.GetPayloads() {
		jobs = append(jobs, &EventAction{dialects.NewEvent(collection, payload), 1})
	}
	return jobs
}

// Returns the collection pairs
func GetTestCollectionPairs(userId uint32, numberOfPayloads int) (*payload.Collection, *payload.Collection, *payload.Collection) {
	collection := GetTestPayloadCollection(userId, numberOfPayloads)
	wrongCollection := GetTestPayloadCollection(userId, numberOfPayloads)
	wrongCollection.Session = proto.String("not-valid-session")
	incompleteCollection := GetTestPayloadCollection(userId, numberOfPayloads)
	for _, payload := range incompleteCollection.GetPayloads() {
		payload.At = nil
	}
	return collection, wrongCollection, incompleteCollection
}

// Returns a valid protobuf collection in bytes and the related jobs
func GetTestProtobufCollectionBody(userId uint32, numberOfPayloads int) (*TrackBodyCollection, []*EventAction) {
	collection, wrongCollection, incompleteCollection := GetTestCollectionPairs(userId, numberOfPayloads)
	collectionBytestream, _ := proto.Marshal(collection)
	disturbedCollectionBytestream, _ := proto.Marshal(wrongCollection)
	incompleteCollectionByteStream, _ := proto.Marshal(incompleteCollection)
	return &TrackBodyCollection{collectionBytestream, disturbedCollectionBytestream, incompleteCollectionByteStream}, GetJobsFromCollection(collection)
}

// Returns a valid json collection in bytes and the related jobs
func GetTestJSONCollectionBody(userId uint32, numberOfPayloads int) (*TrackBodyCollection, []*EventAction) {
	collection, wrongCollection, incompleteCollection := GetTestCollectionPairs(userId, numberOfPayloads)
	var collectionBytestream bytes.Buffer
	var disturbedCollectionBytestream bytes.Buffer
	var incompleteCollectionByteStream bytes.Buffer

	m := jsonpb.Marshaler{false, false, "", true}
	m.Marshal(&collectionBytestream, collection)
	m.Marshal(&disturbedCollectionBytestream, wrongCollection)
	m.Marshal(&incompleteCollectionByteStream, incompleteCollection)

	t := &TrackBodyCollection{
		collectionBytestream.Bytes(),
		disturbedCollectionBytestream.Bytes(),
		incompleteCollectionByteStream.Bytes()}
	return t, GetJobsFromCollection(collection)
}

// Masks all of the jobs
func MaskJobs(jobs []*EventAction) []*EventAction {
	var returnJobs []*EventAction
	for _, j := range jobs {
		event := *j.GetEvent()
		event.TruncateIPv4LastOctet()
		returnJobs = append(returnJobs, &EventAction{&event, j.Attempt})
	}
	return returnJobs
}

// Input cases for the TrackHandler
type TrackHandlerInput struct {
	BodyCollection *TrackBodyCollection
	Time           string
	ContentType    string
	Jobs           []*EventAction
	MaxTestCase    int
}

// Correct and incorrect header generation functions
type HeaderFunction func(t *TrackHandlerInput, fn BodyFunction) map[string]string

func GetValidHeaderWithoutSignature(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"Content-Type": t.ContentType}
}
func GetMissingHeader(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{}
}
func GetHeaderWithoutTime(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time), "Content-Type": t.ContentType}
}
func GetHeaderWithoutSignature(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Time": t.Time, "Content-Type": t.ContentType}
}
func GetHeaderWithInvalidSignature(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Time": t.Time, "X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time) + "x", "Content-Type": t.ContentType}
}
func GetHeaderWithoutContentType(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Time": t.Time, "X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time)}
}
func GetHeaderWithInvalidContentType(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Time": t.Time, "X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time), "Content-Type": "not-existing"}
}
func GetHeaderWithWrongContentType(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	wContentType := map[string]string{"application/json": "application/protobuf", "application/protobuf": "application/json"}[t.ContentType]
	return map[string]string{"X-Hamustro-Time": t.Time, "X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time), "Content-Type": wContentType}
}
func GetValidHeader(t *TrackHandlerInput, fn BodyFunction) map[string]string {
	return map[string]string{"X-Hamustro-Time": t.Time, "X-Hamustro-Signature": GetSignature(fn(t.BodyCollection), t.Time), "Content-Type": t.ContentType}
}

// Signature values
const Optional = 0
const Required = 1
const Any = 2

// Track handler cases
type TrackHandlerTestCase struct {
	Method        string
	GetHeader     HeaderFunction
	GetBody       BodyFunction
	IsTerminating bool
	Signature     int
	ExpectedCode  int
	CheckResults  bool
}

// Executes the test cases for the given inputs
func RunBatchTestOnTrackHandler(t *testing.T, cases []*TrackHandlerTestCase, inputs []*TrackHandlerInput) {
	for i, c := range cases {
		isTerminating = c.IsTerminating
		for _, signature := range map[int][]bool{Optional: []bool{false}, Required: []bool{true}, Any: []bool{true, false}}[c.Signature] {
			signatureRequired = signature
			for _, masked := range []bool{false, true} {
				config.MaskedIP = masked
				for _, isVerbose := range []bool{true, false} {
					verbose = isVerbose         // Sets the verbose mode
					exp = map[string]struct{}{} // Resets the expectations dict

					for j, b := range inputs {
						if b.MaxTestCase <= i {
							continue
						}
						t.Logf("Working on %d/%d test case in %s %s mode %s", i+1, j+1,
							map[bool]string{true: "verbose", false: "production"}[verbose],
							map[bool]string{true: "masked", false: "unmasked"}[masked],
							map[bool]string{true: "with signature", false: "without signature"}[signature])

						// Creates a new request
						req, _ := http.NewRequest(c.Method, "/api/v1/track", bytes.NewBuffer(c.GetBody(b.BodyCollection)))

						// Set up the headers based on the predefined function
						for key, value := range c.GetHeader(b, c.GetBody) {
							req.Header.Set(key, value)
						}
						resp := httptest.NewRecorder()

						// Set up the excepted jobs
						jobs := b.Jobs
						if masked {
							jobs = MaskJobs(jobs)
						}

						// Checks the results and sets the expectations
						if b.Jobs != nil && c.CheckResults {
							for _, job := range jobs {
								SetEventExpectation([]*EventAction{job}, false, false)
							}
						}

						TrackHandler(resp, req) // Calls the API

						// If we're expecting some output, we'll wait for the results
						if b.Jobs != nil && c.CheckResults {
							time.Sleep(150 * time.Millisecond)
							ValidateSending()
						}

						// Log the output to double-check the test case vs reality (debug)
						if verbose && resp.Body.Len() != 0 {
							t.Logf("- Response's body was %s", resp.Body)
						}

						if resp.Code != c.ExpectedCode {
							t.Errorf("Non-expected status code %d with the following body `%s`, it should be %d", resp.Code, resp.Body, c.ExpectedCode)
						}
					}
				}
			}
		}
	}
}

// Tests the API
func TestTrackHandlerRequiredSignature(t *testing.T) {
	t.Log("Creating new workers")
	config = &Config{}                                              // Create an empty config
	storageClient = &SimpleStorageClient{}                          // Define the Simple Storage as a storage
	jobQueue = make(chan Job, 10)                                   // Creates a jobQueue
	log.SetOutput(ioutil.Discard)                                   // Disable the logger
	T, response, catched = t, nil, false                            // Set properties for the SimpleStorageClient
	dispatcher := NewDispatcher(2, &WorkerOptions{RetryAttempt: 5}) // Creates a dispatcher
	dispatcher.Run()                                                // Starts the dispatcher
	config = &Config{SharedSecret: "ultrasafesecret"}               // Creates a config

	if exp := 2; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	var (
		pbNoBody, pbNoBodyJobs                 = GetTestProtobufCollectionBody(98432, 0)
		pbSingleBody, pbSingleBodyJobs         = GetTestProtobufCollectionBody(633289, 1)
		pbMultipleBody, pbMultipleBodyJobs     = GetTestProtobufCollectionBody(53464, 2)
		jsonNoBody, jsonNoBodyJobs             = GetTestJSONCollectionBody(77843, 0)
		jsonSingleBody, jsonSingleBodyJobs     = GetTestJSONCollectionBody(32131, 1)
		jsonMultipleBody, jsonMultipleBodyJobs = GetTestJSONCollectionBody(6546654, 2)
		rTime                                  = "1454514088"
	)

	RunBatchTestOnTrackHandler(t,
		[]*TrackHandlerTestCase{
			{"GET", GetMissingHeader, GetCollectionBody, true, Any, http.StatusServiceUnavailable, false},                        // 1. Service is shutting down
			{"GET", GetMissingHeader, GetCollectionBody, false, Any, http.StatusMethodNotAllowed, false},                         // 2. GET is not supported
			{"POST", GetMissingHeader, GetCollectionBody, false, Required, http.StatusMethodNotAllowed, false},                   // 3. Missing headers with signature
			{"POST", GetMissingHeader, GetCollectionBody, false, Optional, http.StatusBadRequest, false},                         // 4. Missing content-type without signature
			{"POST", GetHeaderWithoutTime, GetCollectionBody, false, Any, http.StatusMethodNotAllowed, false},                    // 5. Missing X-Hamustro-rTime
			{"POST", GetHeaderWithoutSignature, GetCollectionBody, false, Any, http.StatusMethodNotAllowed, false},               // 6. Missing X-Hamustro-Signature
			{"POST", GetHeaderWithoutContentType, GetCollectionBody, false, Any, http.StatusBadRequest, false},                   // 7. Content type is missing
			{"POST", GetHeaderWithInvalidSignature, GetCollectionBody, false, Any, http.StatusMethodNotAllowed, false},           // 8. X-Hamustro-Signature is invalid
			{"POST", GetHeaderWithInvalidContentType, GetCollectionBody, false, Any, http.StatusBadRequest, false},               // 9. Content type is invalid
			{"POST", GetHeaderWithWrongContentType, GetCollectionBody, false, Any, http.StatusBadRequest, false},                 // 10. Content type is not valid for content
			{"POST", GetValidHeader, GetDisturbedCollectionBody, false, Any, http.StatusBadRequest, false},                       // 11. Session is not valid
			{"POST", GetValidHeaderWithoutSignature, GetDisturbedCollectionBody, false, Optional, http.StatusBadRequest, false},  // 12. Session is not valid
			{"POST", GetValidHeader, GetIncompleteCollectionBody, false, Any, http.StatusBadRequest, false},                      // 13. Missing required parameters in the body
			{"POST", GetValidHeaderWithoutSignature, GetIncompleteCollectionBody, false, Optional, http.StatusBadRequest, false}, // 14. Missing required parameters in the body without signature
			{"POST", GetValidHeader, GetCollectionBody, false, Any, http.StatusOK, true},                                         // 15. Valid message
			{"POST", GetValidHeaderWithoutSignature, GetCollectionBody, false, Optional, http.StatusOK, true},                    // 16. Valid message without signature
		},
		[]*TrackHandlerInput{
			{&TrackBodyCollection{[]byte("orange"), nil, nil}, rTime, "", nil, 9},
			{pbSingleBody, rTime, "application/protobuf", pbSingleBodyJobs, 16},
			{pbMultipleBody, rTime, "application/protobuf", pbMultipleBodyJobs, 16},
			{jsonSingleBody, rTime, "application/json", jsonSingleBodyJobs, 16},
			{jsonMultipleBody, rTime, "application/json", jsonMultipleBodyJobs, 16},
		})

	RunBatchTestOnTrackHandler(t,
		[]*TrackHandlerTestCase{
			{"POST", GetValidHeader, GetCollectionBody, false, Any, http.StatusNoContent, false},                      // 1. Valid message without content
			{"POST", GetValidHeaderWithoutSignature, GetCollectionBody, false, Optional, http.StatusNoContent, false}, // 2. Valid message without content without signature
		},
		[]*TrackHandlerInput{
			{jsonNoBody, rTime, "application/json", jsonNoBodyJobs, 2},
			{pbNoBody, rTime, "application/protobuf", pbNoBodyJobs, 2},
		})

	dispatcher.Stop()
}
