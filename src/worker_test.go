package main

import (
	"bytes"
	"fmt"
	"github.com/wunderlist/hamustro/src/dialects"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

// Global variables for testing (hacky)
var T *testing.T
var exp map[string]struct{}
var response error = nil
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

func (c *SimpleStorageClient) Save(msg *bytes.Buffer) error {
	catched = true
	if response != nil {
		return response
	}
	time.Sleep(100 * time.Millisecond)
	T.Logf("Validating received message within the SimpleStorageClient")
	msgString := msg.String()

	var mutex = &sync.Mutex{}
	mutex.Lock()
	if _, ok := exp[msgString]; !ok {
		T.Errorf("Expected message was not %s", msgString)
	} else {
		delete(exp, msgString)
	}
	mutex.Unlock()
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
func (c *BufferedStorageClient) Save(msg *bytes.Buffer) error {
	catched = true
	if response != nil {
		return response
	}
	time.Sleep(100 * time.Millisecond)
	T.Logf("Validating received messages within the BufferedStorageClient")
	msgString := msg.String()
	if _, ok := exp[msgString]; !ok {
		T.Errorf("Expected message was not %s", msgString)
	} else {
		delete(exp, msgString)
	}
	return nil
}

// Worker Id testing with NewWorker
func TestFunctionGetIdAndNewWorker(t *testing.T) {
	jobQueue = make(chan Job, 10)
	storageClient = &SimpleStorageClient{}

	t.Log("Creating a worker with 312 id")
	pool := make(chan *Worker, 1)
	worker := NewWorker(312, &WorkerOptions{}, pool)

	if worker.GetId() != 312 {
		t.Errorf("Expected worker's ID was %d but it was %d instead.", 312, worker.GetId())
	}
}

// Increase Penalty test and buffer size calculation
func TestFunctionIncreasePenaltyAndGetBufferSize(t *testing.T) {
	t.Log("Increasing penalty")
	worker := &Worker{BufferSize: 100, Penalty: 1.0}

	cases := []struct {
		ExpectedBufferSize int
		ExpectedPenalty    float32
		ResetBuffer        bool
	}{
		{100, 1.0, false},
		{150, 1.5, false},
		{225, 2.25, false},
		{100, 1.0, true}}

	for _, c := range cases {
		if c.ResetBuffer {
			worker.ResetBuffer()
		}
		if worker.Penalty != c.ExpectedPenalty {
			t.Errorf("Expected penalty was %f but it was %f instead", c.ExpectedPenalty, worker.Penalty)
		}
		if bs := worker.GetBufferSize(); bs != c.ExpectedBufferSize {
			t.Errorf("Expected buffer size was %d but it was %d instead", c.ExpectedBufferSize, bs)
		}
		worker.IncreasePenalty()
	}
}

// Tests the buffer full condition and adding events to the buffer
func TestFunctionBufferFullAndAddEventToBuffer(t *testing.T) {
	t.Log("Testing worker's buffer functions")
	worker := &Worker{BufferSize: 2, Penalty: 1.0}

	cases := []struct {
		ExpectedBufferLength int
		ExpectedBufferFull   bool
	}{
		{1, false},
		{2, true}}

	for _, c := range cases {
		worker.AddEventToBuffer(GetTestEvent(2158942))
		if length := len(worker.BufferedEvents); length != c.ExpectedBufferLength {
			t.Errorf("Expected buffer length was %d but it was %d instead", c.ExpectedBufferLength, length)
		}
		if isFull := worker.IsBufferFull(); isFull != c.ExpectedBufferFull {
			t.Errorf("Buffer should be %s but it was %s instead.", c.ExpectedBufferFull, isFull)
		}
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

// Sets the actions expectation
func SetEventExpectation(eventActions []*EventAction, fail bool, resetExpectation bool) {
	if fail {
		response = fmt.Errorf("Error was intialized for testing")
	} else {
		response = nil
	}
	if resetExpectation {
		exp = map[string]struct{}{}
	}
	var part *bytes.Buffer
	partStr := ""
	for _, action := range eventActions {
		part, _ = dialects.ConvertJSON(action.Event)
		partStr += part.String()
	}
	exp[partStr] = struct{}{}
}

// Sends a single action to action channel
func SendEventActionToJobChannel(workerPool chan *Worker, worker *Worker, action *EventAction) *Worker {
	expBuffer, _ := dialects.ConvertJSON(action.Event)
	exp[expBuffer.String()] = struct{}{}
	worker.JobChannel <- action
	worker = <-workerPool
	return worker
}

// Validates the previous sending
func ValidateSending() {
	if !catched {
		T.Errorf("Worker didn't catch the expected action")
	}
	catched = false
}

// Sets the expected result, send the message, and validate it.
func SetSendValidate(workerPool chan *Worker, worker *Worker, actions []*EventAction, fail bool, resetExpectation bool) *Worker {
	SetEventExpectation(actions, fail, resetExpectation)
	for _, action := range actions {
		worker = SendEventActionToJobChannel(workerPool, worker, action)
	}
	ValidateSending()
	return worker
}

// Sends a single action to action channel
func SendFlushActionToJobChannel(workerPool chan *Worker, worker *Worker) *Worker {
	lastSave := worker.LastSave
	worker.JobChannel <- &FlushAction{worker.ID}
	worker = <-workerPool
	time.Sleep(200 * time.Millisecond)
	if lastSave == worker.LastSave {
		T.Errorf("Worker last save should not be equivalent, because it flushed the events. The last save was %s", lastSave)
	}
	return worker
}

// Tests the simple storage client (not buffered) with a single worker
func TestSimpleStorageClientWorker(t *testing.T) {
	storageClient = &SimpleStorageClient{} // Define the Simple Storage as a storage
	jobQueue = make(chan Job, 10)          // Creates a jobQueue
	log.SetOutput(ioutil.Discard)          // Disable the logger
	T, response, catched = t, nil, false   // Set properties

	// Create a worker
	t.Log("Creating a single worker")
	pool := make(chan *Worker, 1)
	worker := NewWorker(1, &WorkerOptions{RetryAttempt: 2}, pool)
	worker.Start()

	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(1)
	defer worker.Stop(&wg)

	// Start the test
	worker = <-pool

	t.Log("Creating a single action and send it to the worker")
	worker = SetSendValidate(pool, worker, []*EventAction{&EventAction{GetTestEvent(3423543), 1}}, false, true)

	t.Log("Creating an another single action and send it to the worker")
	worker = SetSendValidate(pool, worker, []*EventAction{&EventAction{GetTestEvent(1321), 1}}, false, true)

	t.Log("Send something that will fail and raise an error")
	action := &EventAction{GetTestEvent(43233), 1}
	worker = SetSendValidate(pool, worker, []*EventAction{action}, true, true)
	if action.Attempt != 2 {
		t.Errorf("Job attempt number should be %d and it was %d instead", 2, action.Attempt)
	}

	t.Log("This failed message must be in the jobQueue, try again.")
	if len(jobQueue) != 1 {
		t.Errorf("worker doesn't contain the previous action")
	}
	worker = SetSendValidate(pool, worker, []*EventAction{(<-jobQueue).(*EventAction)}, false, true)

	t.Log("Send something that will fail and raise an error again")
	worker = SetSendValidate(pool, worker, []*EventAction{&EventAction{GetTestEvent(43254534), 1}}, true, true)

	t.Log("This failed message must be in the jobQueue, but let it fail again.")
	if len(jobQueue) != 1 {
		t.Errorf("jobQueue doesn't contain the previous action")
	}

	action = (<-jobQueue).(*EventAction)
	worker = SetSendValidate(pool, worker, []*EventAction{action}, true, true)
	if action.Attempt != 3 {
		t.Errorf("Job attempt number should be %d and it was %d instead", 3, action.Attempt)
	}
	if len(jobQueue) != 0 {
		t.Errorf("jobQueue have to be empty because it was dropped after the 2nd attempt")
	}

	t.Log("Flush a simple storage")
	worker = SendFlushActionToJobChannel(pool, worker)
}

// Tests the results of the buffered storage worker's test
func CheckResultsForBufferedStorage(worker *Worker, bufferLength int, penalty float32, bufferSize int) {
	if len(worker.BufferedEvents) != bufferLength {
		T.Errorf("Worker's buffered events count should be %d but it was %d instead", bufferLength, len(worker.BufferedEvents))
	}
	if worker.Penalty != penalty {
		T.Errorf("Expected worker's penalty was %d but it was %d instead", penalty, worker.Penalty)
	}
	if worker.GetBufferSize() != bufferSize {
		T.Errorf("Expected worker's buffer size after the error was %d but it was %d instead", bufferSize, worker.GetBufferSize())
	}
}

// Tests the simple storage client (not buffered) with a single worker
func TestBufferedStorageClientWorker(t *testing.T) {
	var wg sync.WaitGroup
	storageClient = &BufferedStorageClient{} // Define the Buffer Storage as a storage
	jobQueue = make(chan Job, 10)            // Creates a jobQueue
	log.SetOutput(ioutil.Discard)            // Disable the logger
	T, response, catched = t, nil, false     // Set properties

	// Create a worker
	t.Log("Creating a single worker with buffer size: 4")
	pool := make(chan *Worker, 1)
	worker := NewWorker(1, &WorkerOptions{BufferSize: 4}, pool)
	worker.Start()

	// Start the test
	worker = <-pool

	t.Log("Creating 4 job and send it to the worker")
	actions := []*EventAction{&EventAction{GetTestEvent(54354353), 1}, &EventAction{GetTestEvent(543), 1}, &EventAction{GetTestEvent(765342), 1}, &EventAction{GetTestEvent(1), 1}}
	SetEventExpectation(actions, false, true)
	for i, action := range actions {
		if len(worker.BufferedEvents) != i {
			t.Errorf("Worker's buffered events count should be %d but it was %d instead", i, len(worker.BufferedEvents))
		}
		worker = SendEventActionToJobChannel(pool, worker, action)
	}
	ValidateSending()
	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

	t.Log("Creating 6 action and send it to the worker, during the process it'll fail after the 4th and will be accepted after the 6th")
	actions = []*EventAction{&EventAction{GetTestEvent(423), 1}, &EventAction{GetTestEvent(654645), 1}, &EventAction{GetTestEvent(123123), 1}, &EventAction{GetTestEvent(16548), 1}}
	SetSendValidate(pool, worker, actions, true, true)
	CheckResultsForBufferedStorage(worker, 4, 1.5, 6)

	actions = append(actions, []*EventAction{&EventAction{GetTestEvent(64562), 1}, &EventAction{GetTestEvent(13127), 1}}...)
	SetEventExpectation(actions, false, true)
	for _, action := range actions[4:] {
		worker = SendEventActionToJobChannel(pool, worker, action)
	}
	ValidateSending()
	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

	t.Log("Creating a single action and send it to the worker that will stay in the buffer until the worker stops")
	action := &EventAction{GetTestEvent(9843211), 1}
	SetEventExpectation([]*EventAction{action}, false, true)
	worker = SendEventActionToJobChannel(pool, worker, action)
	CheckResultsForBufferedStorage(worker, 1, 1.0, 4)

	t.Log("Stop the worker and write out the current buffer")
	wg.Add(1)
	worker.Stop(&wg)
	wg.Wait()
	ValidateSending()
	CheckResultsForBufferedStorage(worker, 0, 1.0, 4)

	// Create a worker again
	t.Log("Creating a single worker again with buffer size: 4")
	pool = make(chan *Worker, 1)
	worker = NewWorker(1, &WorkerOptions{BufferSize: 4}, pool)
	worker.Start()

	// Grab a channel
	worker = <-pool

	t.Log("Creating a single action and send it to the worker that will stay in the buffer until the worker stops")
	action = &EventAction{GetTestEvent(5435), 1}
	SetEventExpectation([]*EventAction{action}, true, true)
	worker = SendEventActionToJobChannel(pool, worker, action)
	CheckResultsForBufferedStorage(worker, 1, 1.0, 4)

	t.Log("Stop the worker and write out the current buffer that will fail and the message will be lost")
	wg.Add(1)
	worker.Stop(&wg)
	wg.Wait()
	ValidateSending()
	CheckResultsForBufferedStorage(worker, 1, 1.5, 6)

	// Create a worker again
	t.Log("Creating a single worker again with buffer size: 4")
	pool = make(chan *Worker, 1)
	worker = NewWorker(1, &WorkerOptions{BufferSize: 4}, pool)
	worker.Start()

	// Grab a channel
	worker = <-pool

	t.Log("Creating 3 action and send it to the worker, during the process the worker gets a flush action")
	actions = []*EventAction{&EventAction{GetTestEvent(4223), 1}, &EventAction{GetTestEvent(66666), 1}, &EventAction{GetTestEvent(969482), 1}}
	SetEventExpectation(actions, false, true)
	for _, action := range actions {
		worker = SendEventActionToJobChannel(pool, worker, action)
	}
	CheckResultsForBufferedStorage(worker, 3, 1, 4)
	worker = SendFlushActionToJobChannel(pool, worker)
	ValidateSending()
	CheckResultsForBufferedStorage(worker, 0, 1, 4)

	t.Log("Stop the worker and write out the current buffer that will fail and the message will be lost")
	wg.Add(1)
	worker.Stop(&wg)
	wg.Wait()
	CheckResultsForBufferedStorage(worker, 0, 1, 4)

	// Create a worker again
	t.Log("Creating a single worker again with buffer size: 4")
	pool = make(chan *Worker, 1)
	worker = NewWorker(1, &WorkerOptions{BufferSize: 3}, pool)
	worker.Start()

	// Grab a channel
	worker = <-pool

	t.Log("Creating 4 action and send it to the worker, during the process the worker gets a flush action")
	actions = []*EventAction{&EventAction{GetTestEvent(7632), 1}, &EventAction{GetTestEvent(3423), 1}, &EventAction{GetTestEvent(23), 1}}
	SetSendValidate(pool, worker, actions, false, true)
	CheckResultsForBufferedStorage(worker, 0, 1, 3)

	actions = append(actions, []*EventAction{&EventAction{GetTestEvent(7532233), 1}}...)
	SetEventExpectation(actions, false, true)
	for _, action := range actions[3:] {
		worker = SendEventActionToJobChannel(pool, worker, action)
	}
	CheckResultsForBufferedStorage(worker, 1, 1.0, 3)

	worker = SendFlushActionToJobChannel(pool, worker)
	CheckResultsForBufferedStorage(worker, 0, 1, 3)
	ValidateSending()

	t.Log("Stop the worker and write out the current buffer that will fail and the message will be lost")
	wg.Add(1)
	worker.Stop(&wg)
	wg.Wait()
	CheckResultsForBufferedStorage(worker, 0, 1, 3)
}

// Waiting for all workers to finish
func WaitingForWorkersToFinish(workerPool chan *Worker, workers []*Worker) {
	var finishedWorker *Worker
	for _, expWorker := range workers {
		for expWorker != finishedWorker {
			finishedWorker = <-workerPool
			expWorker.WorkerPool <- expWorker
		}
	}
}

// Multiple worker tests for simple storage client
func TestSimpleStorageClientMultipleWorker(t *testing.T) {
	t.Log("Testing multiple worker's behaviour")

	// Disable the logger
	log.SetOutput(ioutil.Discard)

	// Define the action Queue and the Buffered Storage Client
	jobQueue = make(chan Job, 10)
	storageClient = &SimpleStorageClient{}

	// Make testing.T and the response global
	T = t
	response = nil
	catched = false

	// Create a worker
	t.Log("Creating two worker to compete with each other")
	pool := make(chan *Worker, 2)
	w1 := NewWorker(1, &WorkerOptions{RetryAttempt: 3}, pool)
	w1.Start()

	w2 := NewWorker(2, &WorkerOptions{RetryAttempt: 3}, pool)
	w2.Start()

	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(2)
	defer w1.Stop(&wg)
	defer w2.Stop(&wg)

	// Create two actions and send it to channels
	action1 := EventAction{GetTestEvent(1262473173), 1}
	expBuffer1, _ := dialects.ConvertJSON(action1.Event)

	action2 := EventAction{GetTestEvent(53484332), 1}
	expBuffer2, _ := dialects.ConvertJSON(action2.Event)

	exp = map[string]struct{}{expBuffer1.String(): {}, expBuffer2.String(): {}}

	// It should catch a different worker with the expected results
	worker := <-pool
	worker.JobChannel <- &action1
	worker = <-pool
	worker.JobChannel <- &action2

	// Get channel to wait until the previous actions are finished
	WaitingForWorkersToFinish(pool, []*Worker{w1, w2})

	if !catched {
		t.Errorf("Worker didn't catch the expected actions")
	}
}

// Multiple worker tests for buffered storage client
func TestBufferedStorageClientMultipleWorker(t *testing.T) {
	t.Log("Testing multiple worker's behaviour for buffered storage")

	storageClient = &BufferedStorageClient{} // Define the Buffer Storage as a storage
	jobQueue = make(chan Job, 10)            // Creates a jobQueue
	log.SetOutput(ioutil.Discard)            // Disable the logger
	T, response, catched = t, nil, false     // Set properties

	// Create workers
	t.Log("Creating two worker with buffer size 2 and 3")

	pool := make(chan *Worker, 2)

	worker1 := NewWorker(1, &WorkerOptions{BufferSize: 2}, pool)
	worker1.Start()

	worker2 := NewWorker(1, &WorkerOptions{BufferSize: 3}, pool)
	worker2.Start()

	// Stop the worker on the end
	var wg sync.WaitGroup
	wg.Add(2)
	defer worker1.Stop(&wg)
	defer worker2.Stop(&wg)

	t.Log("Creating 2 action and send it to the workers")
	// Create two actions and send it to channels
	action1 := EventAction{GetTestEvent(5541289), 1}
	action2 := EventAction{GetTestEvent(7851126), 1}

	// All possible expected result
	expBuffer1, _ := storageClient.GetBatchConverter()([]*dialects.Event{action1.GetEvent()})
	expBuffer2, _ := storageClient.GetBatchConverter()([]*dialects.Event{action2.GetEvent()})
	expBuffer3, _ := storageClient.GetBatchConverter()([]*dialects.Event{action1.GetEvent(), action2.GetEvent()})
	expBuffer4, _ := storageClient.GetBatchConverter()([]*dialects.Event{action2.GetEvent(), action1.GetEvent()})

	exp = map[string]struct{}{expBuffer1.String(): {}, expBuffer2.String(): {}, expBuffer3.String(): {}, expBuffer4.String(): {}}

	// It should catch a different worker with the expected results
	worker := <-pool
	worker.JobChannel <- &action1
	worker = <-pool
	worker.JobChannel <- &action2

	// Get channel to wait until the previous actions are finished
	WaitingForWorkersToFinish(pool, []*Worker{worker1, worker2})
	time.Sleep(100 * time.Millisecond)

	if expLength := 2; len(worker1.BufferedEvents)+len(worker2.BufferedEvents) != expLength {
		T.Errorf("Workers' buffered events count should be %d but it was %d instead", expLength, len(worker1.BufferedEvents)+len(worker2.BufferedEvents))
	}

	t.Log("Flushing the workers")
	worker = SendFlushActionToJobChannel(pool, worker1)
	worker = SendFlushActionToJobChannel(pool, worker2)

	if expLength := 0; len(worker1.BufferedEvents)+len(worker2.BufferedEvents) != expLength {
		T.Errorf("Workers' buffered events count should be %d but it was %d instead", expLength, len(worker1.BufferedEvents)+len(worker2.BufferedEvents))
	}
	ValidateSending()
}
