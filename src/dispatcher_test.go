package main

import (
	"github.com/sub-ninja/hamustro/src/dialects"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

// Testing retry attempt settings
func TestRetryAttemptNewDispatcher(t *testing.T) {
	t.Log("Testing retry attempt settings")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{RetryAttempt: 5}
	dispatcher := NewDispatcher(4, options)
	dispatcher.Start()
	if exp := 4; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}
	for _, w := range dispatcher.Workers {
		if exp := 5; w.RetryAttempt != exp {
			t.Errorf("Expected %d worker's retry attempt property was %d and it was %d instead", w.ID, exp, w.RetryAttempt)
		}
	}
}

// Testing buffer size calculation for not spreading buffers
func TestNotSpreadBufferNewDispatcher(t *testing.T) {
	t.Log("Testing buffer size calculation for not spreading buffers")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{SpreadBuffer: false, BufferSize: 10000}
	dispatcher := NewDispatcher(4, options)
	dispatcher.Start()
	if exp := 4; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}
	for _, w := range dispatcher.Workers {
		if exp := 10000; w.BufferSize != exp {
			t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", w.ID, exp, w.BufferSize)
		}
	}
}

// Testing buffer size calculation for not spreading buffers
func TestSpreadBufferNewDispatcher(t *testing.T) {
	t.Log("Testing buffer size calculation for not spreading buffers")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{SpreadBuffer: true, BufferSize: 10000}
	dispatcher := NewDispatcher(3, options)
	dispatcher.Start()
	if exp := 3; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	w := dispatcher.Workers[0]
	if exp := 7500; w.BufferSize != exp {
		t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", w.ID, exp, w.BufferSize)
	}

	w = dispatcher.Workers[1]
	if exp := 10000; w.BufferSize != exp {
		t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", w.ID, exp, w.BufferSize)
	}

	w = dispatcher.Workers[2]
	if exp := 12500; w.BufferSize != exp {
		t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", w.ID, exp, w.BufferSize)
	}
}

// Testing the buffer size caluculation function
func TestFunctionDispatcherGetBufferSize(t *testing.T) {
	t.Log("Testing the buffer size caluculation function")
	dispatcher := &Dispatcher{MaxWorkers: 3, WorkerOptions: &WorkerOptions{BufferSize: 10000, SpreadBuffer: true}}
	if exp := 7500; dispatcher.GetBufferSize(0) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(0))
	}
	if exp := 10000; dispatcher.GetBufferSize(1) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(1))
	}
	if exp := 12500; dispatcher.GetBufferSize(2) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(2))
	}

	dispatcher = &Dispatcher{MaxWorkers: 3, WorkerOptions: &WorkerOptions{BufferSize: 10000, SpreadBuffer: false}}
	if exp := 10000; dispatcher.GetBufferSize(0) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(0))
	}
	if exp := 10000; dispatcher.GetBufferSize(1) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(1))
	}
	if exp := 10000; dispatcher.GetBufferSize(2) != exp {
		t.Errorf("Expected buffer size was %d and it was %d instead", exp, dispatcher.GetBufferSize(2))
	}
}

// Testing the dispatcher listen function
func TestDispatcherListen(t *testing.T) {
	t.Log("Testing the dispatcher listen function")

	// Define the job Queue and the Buffered Storage Client
	storageClient = &SimpleStorageClient{}
	jobQueue = make(chan *Job, 10)

	// Disable the logger
	log.SetOutput(ioutil.Discard)

	// Testing responses
	T = t
	response = nil
	catched = false

	// Creates the dispatcher and listen for new jobs
	options := &WorkerOptions{RetryAttempt: 5}
	dispatcher := NewDispatcher(2, options)
	dispatcher.Run()

	if exp := 2; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	// Create two jobs and put it into the job queue
	t.Log("Creating two jobs and put it into the job queue")
	job1 := Job{GetTestEvent(423432), 1}
	expBuffer1, _ := dialects.ConvertJSON(job1.Event)

	job2 := Job{GetTestEvent(7643329), 1}
	expBuffer2, _ := dialects.ConvertJSON(job2.Event)

	exp = map[string]struct{}{expBuffer1.String(): {}, expBuffer2.String(): {}}

	// It should catch a different worker with the expected results
	jobQueue <- &job1
	jobQueue <- &job2

	// Wait until both is finished
	time.Sleep(150 * time.Millisecond)

	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}

	t.Log("Tries to stop the workers")
	// Stops the workers
	dispatcher.Stop()
}
