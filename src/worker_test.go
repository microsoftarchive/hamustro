package main

import (
	// 	"bytes"
	// 	"fmt"
	"github.com/wunderlist/hamustro/src/dialects"
	// 	"io/ioutil"
	// 	"log"
	// 	"sync"
	// 	"testing"
	// 	"time"
)

// // Global variables for testing (hacky)
// var T *testing.T
// var exp map[string]struct{}
// var response error = nil
// var catched bool = false

// // Simple (not buffered) Storage Client for testing
// type SimpleStorageClient struct{}

// func (c *SimpleStorageClient) IsBufferedStorage() bool {
// 	return false
// }
// func (c *SimpleStorageClient) GetConverter() dialects.Converter {
// 	return dialects.ConvertJSON
// }
// func (c *SimpleStorageClient) GetBatchConverter() dialects.BatchConverter {
// 	return nil
// }
// func (c *SimpleStorageClient) Save(msg *bytes.Buffer) error {
// 	catched = true
// 	if response != nil {
// 		return response
// 	}
// 	time.Sleep(100 * time.Millisecond)
// 	T.Logf("Validating received message within the SimpleStorageClient")
// 	msgString := msg.String()
// 	if _, ok := exp[msgString]; !ok {
// 		T.Errorf("Expected message was not %s", msgString)
// 	} else {
// 		delete(exp, msgString)
// 	}
// 	return nil
// }

// // Buffered Storage Client for testing
// type BufferedStorageClient struct{}

// func (c *BufferedStorageClient) IsBufferedStorage() bool {
// 	return true
// }
// func (c *BufferedStorageClient) GetConverter() dialects.Converter {
// 	return nil
// }
// func (c *BufferedStorageClient) GetBatchConverter() dialects.BatchConverter {
// 	return dialects.ConvertBatchJSON
// }
// func (c *BufferedStorageClient) Save(msg *bytes.Buffer) error {
// 	catched = true
// 	if response != nil {
// 		return response
// 	}
// 	time.Sleep(100 * time.Millisecond)
// 	T.Logf("Validating received messages within the BufferedStorageClient")
// 	msgString := msg.String()
// 	if _, ok := exp[msgString]; !ok {
// 		T.Errorf("Expected message was not %s", msgString)
// 	} else {
// 		delete(exp, msgString)
// 	}
// 	return nil
// }

// // Worker Id testing with NewWorker
// func TestFunctionGetIdAndNewWorker(t *testing.T) {
// 	jobQueue = make(chan *Job, 10)
// 	storageClient = &SimpleStorageClient{}

// 	t.Log("Creating a worker with 312 id")
// 	pool := make(chan chan *Job, 1)
// 	worker := NewWorker(312, &WorkerOptions{}, pool)

// 	if worker.GetId() != 312 {
// 		t.Errorf("Expected worker's ID was %d but it was %d instead.", 312, worker.GetId())
// 	}
// }

// // Increase Penalty test and buffer size calculation
// func TestFunctionIncreasePenaltyAndGetBufferSize(t *testing.T) {
// 	t.Log("Increasing penalty")
// 	worker := &Worker{BufferSize: 100, Penalty: 1.0}

// 	cases := []struct {
// 		ExpectedBufferSize int
// 		ExpectedPenalty    float32
// 		ResetBuffer        bool
// 	}{
// 		{100, 1.0, false},
// 		{150, 1.5, false},
// 		{225, 2.25, false},
// 		{100, 1.0, true}}

// 	for _, c := range cases {
// 		if c.ResetBuffer {
// 			worker.ResetBuffer()
// 		}
// 		if worker.Penalty != c.ExpectedPenalty {
// 			t.Errorf("Expected penalty was %f but it was %f instead", c.ExpectedPenalty, worker.Penalty)
// 		}
// 		if bs := worker.GetBufferSize(); bs != c.ExpectedBufferSize {
// 			t.Errorf("Expected buffer size was %d but it was %d instead", c.ExpectedBufferSize, bs)
// 		}
// 		worker.IncreasePenalty()
// 	}
// }

// // Tests the buffer full condition and adding events to the buffer
// func TestFunctionBufferFullAndAddEventToBuffer(t *testing.T) {
// 	t.Log("Testing worker's buffer functions")
// 	worker := &Worker{BufferSize: 2, Penalty: 1.0}

// 	cases := []struct {
// 		ExpectedBufferLength int
// 		ExpectedBufferFull   bool
// 	}{
// 		{1, false},
// 		{2, true}}

// 	for _, c := range cases {
// 		worker.AddEventToBuffer(GetTestEvent(2158942))
// 		if length := len(worker.BufferedEvents); length != c.ExpectedBufferLength {
// 			t.Errorf("Expected buffer length was %d but it was %d instead", c.ExpectedBufferLength, length)
// 		}
// 		if isFull := worker.IsBufferFull(); isFull != c.ExpectedBufferFull {
// 			t.Errorf("Buffer should be %s but it was %s instead.", c.ExpectedBufferFull, isFull)
// 		}
// 	}
// 	worker.ResetBuffer()
// 	if exp := 0; len(worker.BufferedEvents) != exp {
// 		t.Errorf("Expected buffer length was %d but it was %d instead", exp, len(worker.BufferedEvents))
// 	}
// }

// Returns an Event for testing purposes
func GetTestEvent(userId uint32) *dialects.Event {
	return &dialects.Event{
		DeviceID:       "a73b1c37-2c24-4786-af7a-16de88fbe23a",
		ClientID:       "bce44f67b2661fd445d469b525b04f68",
		Session:        "244f056dee6d475ec673ea0d20b69bab",
		Nr:             1,
		SystemVersion:  "10.10",
		ProductVersion: "1.1.2",
		At:             "2016-02-05T15:05:04",
		Event:          "Client.CreateUser",
		System:         "OSX",
		ProductGitHash: "5416a5889392d509e3bafcf40f6388e83aab23e6",
		UserID:         userId,
		IP:             "214.160.227.22",
		Parameters:     "",
		IsTesting:      false}
}

// // Sets the jobs expectation
// func SetJobExpectation(jobs []*Job, fail bool, resetExpectation bool) {
// 	if fail {
// 		response = fmt.Errorf("Error was intialized for testing")
// 	} else {
// 		response = nil
// 	}
// 	if resetExpectation {
// 		exp = map[string]struct{}{}
// 	}
// 	var part *bytes.Buffer
// 	partStr := ""
// 	for _, job := range jobs {
// 		part, _ = dialects.ConvertJSON(job.Event)
// 		partStr += part.String()
// 	}
// 	exp[partStr] = struct{}{}
// }

// // Sends a single job to job channel
// func SendToJobChannel(workerPool chan chan *Job, jobChannel chan *Job, job *Job) chan *Job {
// 	expBuffer, _ := dialects.ConvertJSON(job.Event)
// 	exp[expBuffer.String()] = struct{}{}
// 	jobChannel <- job
// 	jobChannel = <-workerPool
// 	return jobChannel
// }

// // Validates the previous sending
// func ValidateSending() {
// 	if !catched {
// 		T.Errorf("Worker didn't catch the expected job")
// 	}
// 	catched = false
// }

// // Sets the expected result, send the message, and validate it.
// func SetSendValidate(workerPool chan chan *Job, jobChannel chan *Job, jobs []*Job, fail bool, resetExpectation bool) chan *Job {
// 	SetJobExpectation(jobs, fail, resetExpectation)
// 	for _, job := range jobs {
// 		jobChannel = SendToJobChannel(workerPool, jobChannel, job)
// 	}
// 	ValidateSending()
// 	return jobChannel
// }

// // Tests the simple storage client (not buffered) with a single worker
// func TestSimpleStorageClientWorker(t *testing.T) {
// 	storageClient = &SimpleStorageClient{} // Define the Simple Storage as a storage
// 	jobQueue = make(chan *Job, 10)         // Creates a jobQueue
// 	log.SetOutput(ioutil.Discard)          // Disable the logger
// 	T, response, catched = t, nil, false   // Set properties

// 	// Create a worker
// 	t.Log("Creating a single worker")
// 	pool := make(chan chan *Job, 1)
// 	worker := NewWorker(1, &WorkerOptions{RetryAttempt: 2}, pool)
// 	worker.Start()

// 	// Stop the worker on the end
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	defer worker.Stop(&wg)

// 	// Start the test
// 	jobChannel := <-pool

// 	t.Log("Creating a single job and send it to the worker")
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{&Job{GetTestEvent(3423543), 1}}, false, true)

// 	t.Log("Creating an another single job and send it to the worker")
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{&Job{GetTestEvent(1321), 1}}, false, true)

// 	t.Log("Send something that will fail and raise an error")
// 	job := &Job{GetTestEvent(43233), 1}
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{job}, true, true)
// 	if job.Attempt != 2 {
// 		t.Errorf("Job attempt number should be %d and it was %d instead", 2, job.Attempt)
// 	}

// 	t.Log("This failed message must be in the jobQueue, try again.")
// 	if len(jobQueue) != 1 {
// 		t.Errorf("jobChannel doesn't contain the previous job")
// 	}
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{<-jobQueue}, false, true)

// 	t.Log("Send something that will fail and raise an error again")
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{&Job{GetTestEvent(43254534), 1}}, true, true)

// 	t.Log("This failed message must be in the jobQueue, but let it fail again.")
// 	if len(jobQueue) != 1 {
// 		t.Errorf("jobQueue doesn't contain the previous job")
// 	}
// 	job = <-jobQueue
// 	jobChannel = SetSendValidate(pool, jobChannel, []*Job{job}, true, true)
// 	if job.Attempt != 3 {
// 		t.Errorf("Job attempt number should be %d and it was %d instead", 3, job.Attempt)
// 	}
// 	if len(jobQueue) != 0 {
// 		t.Errorf("jobQueue have to be empty because it was dropped after the 2nd attempt")
// 	}
// }

// // Tests the results of the buffered storage worker's test
// func CheckResultsForBufferedStorage(worker *Worker, bufferLength int, penalty float32, bufferSize int) {
// 	if len(worker.BufferedEvents) != bufferLength {
// 		T.Errorf("Worker's buffered events count should be %d but it was %d instead", bufferLength, len(worker.BufferedEvents))
// 	}
// 	if worker.Penalty != penalty {
// 		T.Errorf("Expected worker's penalty was %d but it was %d instead", penalty, worker.Penalty)
// 	}
// 	if worker.GetBufferSize() != bufferSize {
// 		T.Errorf("Expected worker's buffer size after the error was %d but it was %d instead", bufferSize, worker.GetBufferSize())
// 	}
// }

// // Tests the simple storage client (not buffered) with a single worker
// func TestBufferedStorageClientWorker(t *testing.T) {
// 	var wg sync.WaitGroup
// 	storageClient = &BufferedStorageClient{} // Define the Buffer Storage as a storage
// 	jobQueue = make(chan *Job, 10)           // Creates a jobQueue
// 	log.SetOutput(ioutil.Discard)            // Disable the logger
// 	T, response, catched = t, nil, false     // Set properties

// 	// Create a worker
// 	t.Log("Creating a single worker with buffer size: 4")
// 	pool := make(chan chan *Job, 1)
// 	worker := NewWorker(1, &WorkerOptions{BufferSize: 4}, pool)
// 	worker.Start()

// 	// Start the test
// 	jobChannel := <-pool

// 	t.Log("Creating 3 job and send it to the worker")
// 	jobs := []*Job{&Job{GetTestEvent(54354353), 1}, &Job{GetTestEvent(543), 1}, &Job{GetTestEvent(765342), 1}, &Job{GetTestEvent(1), 1}}
// 	SetJobExpectation(jobs, false, true)
// 	for i, job := range jobs {
// 		if len(worker.BufferedEvents) != i {
// 			t.Errorf("Worker's buffered events count should be %d but it was %d instead", i, len(worker.BufferedEvents))
// 		}
// 		jobChannel = SendToJobChannel(pool, jobChannel, job)
// 	}
// 	ValidateSending()
// 	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

// 	t.Log("Creating 6 job and send it to the worker, during the process it'll fail after the 4th and will be accepted after the 6th")
// 	jobs = []*Job{&Job{GetTestEvent(423), 1}, &Job{GetTestEvent(654645), 1}, &Job{GetTestEvent(123123), 1}, &Job{GetTestEvent(16548), 1}}
// 	SetSendValidate(pool, jobChannel, jobs, true, true)
// 	CheckResultsForBufferedStorage(worker, 4, 1.5, 6)

// 	jobs = append(jobs, []*Job{&Job{GetTestEvent(64562), 1}, &Job{GetTestEvent(13127), 1}}...)
// 	SetJobExpectation(jobs, false, true)
// 	for _, job := range jobs[4:] {
// 		jobChannel = SendToJobChannel(pool, jobChannel, job)
// 	}
// 	ValidateSending()
// 	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

// 	t.Log("Creating a single job and send it to the worker that will stay in the buffer until the worker stops")
// 	job := &Job{GetTestEvent(9843211), 1}
// 	SetJobExpectation([]*Job{job}, false, true)
// 	jobChannel = SendToJobChannel(pool, jobChannel, job)
// 	CheckResultsForBufferedStorage(worker, 1, 1.0, 4)

// 	t.Log("Stop the worker and write out the current buffer")
// 	wg.Add(1)
// 	worker.Stop(&wg)
// 	wg.Wait()
// 	ValidateSending()
// 	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

// 	// Create a worker again
// 	t.Log("Creating a single worker again with buffer size: 4")
// 	pool = make(chan chan *Job, 1)
// 	worker = NewWorker(1, &WorkerOptions{BufferSize: 4}, pool)
// 	worker.Start()

// 	// Grab a channel
// 	jobChannel = <-pool

// 	t.Log("Creating a single job and send it to the worker that will stay in the buffer until the worker stops")
// 	job = &Job{GetTestEvent(5435), 1}
// 	SetJobExpectation([]*Job{job}, true, true)
// 	jobChannel = SendToJobChannel(pool, jobChannel, job)
// 	CheckResultsForBufferedStorage(worker, 1, 1.0, 4)

// 	t.Log("Stop the worker and write out the current buffer that will fail and the message will be lost")
// 	wg.Add(1)
// 	worker.Stop(&wg)
// 	wg.Wait()
// 	ValidateSending()
// 	CheckResultsForBufferedStorage(worker, 1, 1.0, 4)
// }

// // Multiple worker tests
// func TestMultipleWorker(t *testing.T) {
// 	t.Log("Testing multiple worker's behaviour")

// 	// Disable the logger
// 	log.SetOutput(ioutil.Discard)

// 	// Define the job Queue and the Buffered Storage Client
// 	jobQueue = make(chan *Job, 10)
// 	storageClient = &SimpleStorageClient{}

// 	// Make testing.T and the response global
// 	T = t
// 	response = nil
// 	catched = false

// 	// Create a worker
// 	t.Log("Creating two worker to compete with each other")
// 	pool := make(chan chan *Job, 2)
// 	w1 := NewWorker(1, &WorkerOptions{RetryAttempt: 3}, pool)
// 	w1.Start()

// 	w2 := NewWorker(2, &WorkerOptions{RetryAttempt: 3}, pool)
// 	w2.Start()

// 	// Stop the worker on the end
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	defer w1.Stop(&wg)
// 	defer w2.Stop(&wg)

// 	var jobChannel chan *Job

// 	// Create two jobs and send it to channels
// 	job1 := Job{GetTestEvent(1262473173), 1}
// 	expBuffer1, _ := dialects.ConvertJSON(job1.Event)

// 	job2 := Job{GetTestEvent(53484332), 1}
// 	expBuffer2, _ := dialects.ConvertJSON(job2.Event)

// 	exp = map[string]struct{}{expBuffer1.String(): {}, expBuffer2.String(): {}}

// 	// It should catch a different worker with the expected results
// 	jobChannel = <-pool
// 	jobChannel <- &job1
// 	jobChannel = <-pool
// 	jobChannel <- &job2

// 	// Get two new channel to wait until the previous jobs are finished
// 	<-pool
// 	<-pool

// 	if !catched {
// 		t.Errorf("Worker didn't catch the expected jobs")
// 	}
// }
