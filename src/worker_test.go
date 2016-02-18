package main

import (
	"bytes"
	"fmt"
	"github.com/sub-ninja/hamustro/src/dialects"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

// Global variables for testing (hacky)
var T *testing.T
var exp map[int]*bytes.Buffer
var sResp error = nil
var catched bool = false

// Simple (not buffered) Storage Client for testing
type SimpleStorageClient struct{}

func (c *SimpleStorageClient) IsBufferedStorage() bool {
	return false
}
func (c *SimpleStorageClient) GetConverter() dialects.Converter {
	return dialects.ConvertJSON
}
func (c *SimpleStorageClient) GetBatchConverter() dialects.BatchConverter {
	return nil
}
func (c *SimpleStorageClient) Save(workerID int, msg *bytes.Buffer) error {
	catched = true
	if sResp != nil {
		return sResp
	}
	time.Sleep(100 * time.Millisecond)
	T.Logf("(%d worker) Validating received message within the SimpleStorageClient", workerID)
	if exp[workerID].String() != msg.String() {
		T.Errorf("(%d worker) Expected message was `%s` and it was `%s` instead.", workerID, exp[workerID], msg)
	}
	return nil
}

// Buffered Storage Client for testing
type BufferedStorageClient struct{}

func (c *BufferedStorageClient) IsBufferedStorage() bool {
	return true
}
func (c *BufferedStorageClient) GetConverter() dialects.Converter {
	return nil
}
func (c *BufferedStorageClient) GetBatchConverter() dialects.BatchConverter {
	return dialects.ConvertBatchJSON
}
func (c *BufferedStorageClient) Save(workerID int, msg *bytes.Buffer) error {
	catched = true
	if sResp != nil {
		return sResp
	}
	time.Sleep(100 * time.Millisecond)
	T.Logf("(%d worker) Validating received messages within the BufferedStorageClient", workerID)
	if exp[workerID].String() != msg.String() {
		T.Errorf("(%d worker) Expected message was `%s` and it was `%s` instead.", workerID, exp[workerID], msg)
	}
	return nil
}

// Worker Id testing with NewWorker
func TestFunctionGetIdAndNewWorker(t *testing.T) {
	jobQueue = make(chan *Job, 10)
	storageClient = &SimpleStorageClient{}

	t.Log("Creating a worker with 312 id")
	pool := make(chan chan *Job, 1)
	worker := NewWorker(312, &WorkerOptions{}, pool)

	if worker.GetId() != 312 {
		t.Errorf("Expected worker's ID was %d but it was %d instead.", 312, worker.GetId())
	}
}

// Increase Penalty test and buffer size calculation
func TestFunctionIncreasePenaltyAndGetBufferSize(t *testing.T) {
	t.Log("Increasing penalty")
	worker := &Worker{BufferSize: 100, Penalty: 1.0}
	if exbs := 100; worker.GetBufferSize() != exbs {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, worker.GetBufferSize())
	}
	worker.IncreasePenalty()
	if exp := float32(1.5); worker.Penalty != exp {
		t.Errorf("Expected penalty was %f but it was %f instead", exp, worker.Penalty)
	}
	if exbs := 150; worker.GetBufferSize() != exbs {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, worker.GetBufferSize())
	}
	worker.IncreasePenalty()
	if exp := float32(2.25); worker.Penalty != exp {
		t.Errorf("Expected penalty was %f but it was %f instead", exp, worker.Penalty)
	}
	if exbs := 225; worker.GetBufferSize() != exbs {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, worker.GetBufferSize())
	}
	t.Log("Reset the buffer")
	worker.ResetBuffer()
	if exp := float32(1.0); worker.Penalty != exp {
		t.Errorf("Expected penalty was %f but it was %f instead", exp, worker.Penalty)
	}
	if exbs := 100; worker.GetBufferSize() != exbs {
		t.Errorf("Expected buffer size was %d but it was %d instead", exp, worker.GetBufferSize())
	}
}

// Tests the buffer full condition and adding events to the buffer
func TestFunctionBufferFullAndAddEventToBuffer(t *testing.T) {
	t.Log("Testing worker's buffer functions")
	worker := &Worker{BufferSize: 2, Penalty: 1.0}
	worker.AddEventToBuffer(GetTestEvent(2158942))
	if exp := 1; len(worker.BufferedEvents) != exp {
		t.Errorf("Expected buffer length was %d but it was %d instead", exp, len(worker.BufferedEvents))
	}
	if worker.IsBufferFull() {
		t.Errorf("Buffer should not be full with %d length", len(worker.BufferedEvents))
	}
	worker.AddEventToBuffer(GetTestEvent(54389423))
	if exp := 2; len(worker.BufferedEvents) != exp {
		t.Errorf("Expected buffer length was %d but it was %d instead", exp, len(worker.BufferedEvents))
	}
	if !worker.IsBufferFull() {
		t.Errorf("Buffer should be full with %d length", len(worker.BufferedEvents))
	}
	worker.ResetBuffer()
	if exp := 0; len(worker.BufferedEvents) != exp {
		t.Errorf("Expected buffer length was %d but it was %d instead", exp, len(worker.BufferedEvents))
	}
}

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

// Tests the simple storage client (not buffered) with a single worker
func TestSimpleStorageClientWorker(t *testing.T) {
	exp = make(map[int]*bytes.Buffer)

	// Disable the logger
	log.SetOutput(ioutil.Discard)

	// Define the job Queue and the Simple Storage Client
	jobQueue = make(chan *Job, 10)
	storageClient = &SimpleStorageClient{}

	// Make testing.T and the response global
	T = t
	sResp = nil
	catched = false

	// Create a worker
	t.Log("Creating a single worker")
	pool := make(chan chan *Job, 1)
	worker := NewWorker(1, &WorkerOptions{RetryAttempt: 2}, pool)
	worker.Start()

	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(1)
	defer worker.Stop(&wg)

	// Start the test
	jobChannel := <-pool

	t.Log("Creating a single job and send it to the worker")
	job := Job{GetTestEvent(3423543), 1}
	expBuffer, _ := dialects.ConvertJSON(job.Event)
	exp[worker.ID] = expBuffer
	jobChannel <- &job
	jobChannel = <-pool

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("Creating an another single job and send it to the worker")
	job = Job{GetTestEvent(1321), 1}
	expBuffer, _ = dialects.ConvertJSON(job.Event)
	exp[worker.ID] = expBuffer
	jobChannel <- &job
	jobChannel = <-pool

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("Send something that will fail and raise an error")
	sResp = fmt.Errorf("Error was intialized for testing")
	job = Job{GetTestEvent(43233), 1}
	expBuffer, _ = dialects.ConvertJSON(job.Event)
	exp[worker.ID] = expBuffer
	jobChannel <- &job
	jobChannel = <-pool

	if job.Attempt != 2 {
		t.Errorf("Job attempt number should be %d and it was %d instead", 2, job.Attempt)
	}
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("This failed message must be in the jobQueue, try again.")
	if len(jobQueue) != 1 {
		t.Errorf("jobChannel doesn't contain the previous job")
	}
	jobq := <-jobQueue
	sResp = nil
	jobChannel <- jobq
	jobChannel = <-pool

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("Send something that will fail and raise an error again")
	sResp = fmt.Errorf("Error was intialized for testing")
	job = Job{GetTestEvent(43254534), 1}
	expBuffer, _ = dialects.ConvertJSON(job.Event)
	exp[worker.ID] = expBuffer
	jobChannel <- &job
	jobChannel = <-pool

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("This failed message must be in the jobQueue, but let it fail again.")
	if len(jobQueue) != 1 {
		t.Errorf("jobQueue doesn't contain the previous job")
	}
	jobq = <-jobQueue
	jobChannel <- jobq
	jobChannel = <-pool

	if job.Attempt != 3 {
		t.Errorf("Job attempt number should be %d and it was %d instead", 3, job.Attempt)
	}
	if len(jobQueue) != 0 {
		t.Errorf("jobQueue have to be empty because it was dropped after the 2nd attempt")
	}
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
}

// Tests the simple storage client (not buffered) with a single worker
func TestBufferedStorageClientWorker(t *testing.T) {
	exp = make(map[int]*bytes.Buffer)

	// Disable the logger
	log.SetOutput(ioutil.Discard)

	// Define the job Queue and the Buffered Storage Client
	jobQueue = make(chan *Job, 10)
	storageClient = &BufferedStorageClient{}

	// Make testing.T and the response global
	T = t
	sResp = nil
	catched = false

	// Create a worker
	t.Log("Creating a single worker")
	pool := make(chan chan *Job, 1)
	worker := NewWorker(1, &WorkerOptions{BufferSize: 10}, pool)
	worker.Start()

	// Start the test
	jobChannel := <-pool
	var job Job

	t.Log("Creating 9 job and send it to the worker")
	partStr := ""
	for i := 0; i < 9; i++ {
		job = Job{GetTestEvent(uint32(56746535 + i)), 1}
		part, _ := dialects.ConvertJSON(job.Event)
		partStr += part.String()
		jobChannel <- &job
		jobChannel = <-pool

		if exnr := i + 1; len(worker.BufferedEvents) != exnr {
			t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
		}
	}

	t.Log("Creating the 10th job and send it to the worker that will proceed the buffer")
	job = Job{GetTestEvent(1), 1}
	part, _ := dialects.ConvertJSON(job.Event)
	partStr += part.String()
	exp[worker.ID] = bytes.NewBuffer([]byte(partStr))
	jobChannel <- &job
	jobChannel = <-pool

	if exnr := 0; len(worker.BufferedEvents) != exnr {
		t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
	}
	if expen := float32(1.0); worker.Penalty != expen {
		t.Errorf("Expected worker's penalty was %d but it was %d instead", expen, worker.Penalty)
	}
	if exnr := 10; worker.GetBufferSize() != exnr {
		t.Errorf("Expected worker's buffer size after the error was %d but it was %d instead", exnr, worker.GetBufferSize())
	}
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("Creating 14 job and send it to the worker, during the process it'll fail after the 10th")
	sResp = fmt.Errorf("Error was intialized for testing")
	partStr = ""
	for i := 0; i < 14; i++ {
		job = Job{GetTestEvent(uint32(213432 + i)), 1}
		part, _ := dialects.ConvertJSON(job.Event)
		partStr += part.String()
		jobChannel <- &job
		jobChannel = <-pool

		if exnr := i + 1; len(worker.BufferedEvents) != exnr {
			t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
		}
	}

	if expen := float32(1.5); worker.Penalty != expen {
		t.Errorf("Expected worker's penalty was %d but it was %d instead", expen, worker.Penalty)
	}
	if exnr := 15; worker.GetBufferSize() != exnr {
		t.Errorf("Expected worker's buffer size after the error was %d but it was %d instead", exnr, worker.GetBufferSize())
	}

	sResp = nil
	t.Log("Creating the 15th job and send it to the worker that will proceed the buffer")
	job = Job{GetTestEvent(1), 1}
	part, _ = dialects.ConvertJSON(job.Event)
	partStr += part.String()
	exp[worker.ID] = bytes.NewBuffer([]byte(partStr))
	jobChannel <- &job
	jobChannel = <-pool

	if exnr := 0; len(worker.BufferedEvents) != exnr {
		t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
	}
	if expen := float32(1.0); worker.Penalty != expen {
		t.Errorf("Expected worker's penalty was %d but it was %d instead", expen, worker.Penalty)
	}
	if exnr := 10; worker.GetBufferSize() != exnr {
		t.Errorf("Expected worker's buffer size after the error was %d but it was %d instead", exnr, worker.GetBufferSize())
	}
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
	catched = false

	t.Log("Creating a single job and send it to the worker that will stay in the buffer until the worker stops")
	job = Job{GetTestEvent(1), 1}
	expBuffer, _ := dialects.ConvertJSON(job.Event)
	exp[worker.ID] = expBuffer
	jobChannel <- &job
	jobChannel = <-pool

	if exnr := 1; len(worker.BufferedEvents) != exnr {
		t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
	}

	t.Log("Stop the worker and write out the current buffer")
	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(1)
	worker.Stop(&wg)
	wg.Wait()

	if exnr := 0; len(worker.BufferedEvents) != exnr {
		t.Errorf("Worker's buffered events count should be %d but it was %d instead", exnr, len(worker.BufferedEvents))
	}
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
}

// Multiple worker tests
func TestMultipleWorker(t *testing.T) {
	exp = make(map[int]*bytes.Buffer)
	t.Log("Testing multiple worker's behaviour")

	// Disable the logger
	log.SetOutput(ioutil.Discard)

	// Define the job Queue and the Buffered Storage Client
	jobQueue = make(chan *Job, 10)
	storageClient = &SimpleStorageClient{}

	// Make testing.T and the response global
	T = t
	sResp = nil
	catched = false

	// Create a worker
	t.Log("Creating two worker to compete with each other")
	pool := make(chan chan *Job, 2)
	w1 := NewWorker(1, &WorkerOptions{RetryAttempt: 3}, pool)
	w1.Start()

	w2 := NewWorker(2, &WorkerOptions{RetryAttempt: 3}, pool)
	w2.Start()

	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(2)
	defer w1.Stop(&wg)
	defer w2.Stop(&wg)

	var jobChannel chan *Job

	// Create two jobs and send it to channels
	job1 := Job{GetTestEvent(1262473173), 1}
	expBuffer1, _ := dialects.ConvertJSON(job1.Event)
	exp[w1.ID] = expBuffer1

	job2 := Job{GetTestEvent(53484332), 1}
	expBuffer2, _ := dialects.ConvertJSON(job2.Event)
	exp[w2.ID] = expBuffer2

	// It should catch a different worker with the expected results
	jobChannel = <-pool
	jobChannel <- &job1
	jobChannel = <-pool
	jobChannel <- &job2

	// Get two new channel to wait until the previous jobs are finished
	<-pool
	<-pool

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}
}
