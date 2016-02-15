package main

import (
	"./dialects"
	"bytes"
	"fmt"
	"sync"
	"testing"
)

// Global variables for testing (hacky)
var T *testing.T
var exp *bytes.Buffer
var sResp error = nil

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
	if sResp != nil {
		return sResp
	}
	T.Log("Validating received message within the SimpleStorageClient")
	if exp.String() != msg.String() {
		T.Errorf("Expected message was `%s` and it was `%s` instead.", exp, msg)
	}
	return nil
}

// Tests the simple storage client (not buffered) with a single worker
func TestSimpleStorageClientWorker(t *testing.T) {
	jobQueue = make(chan *Job, 10)
	storageClient = &SimpleStorageClient{}
	T = t
	sResp = nil

	t.Log("Creating a single worker")
	pool := make(chan chan *Job, 2)
	worker := NewWorker(1, 10, pool)
	worker.RetryAttempt = 2
	worker.Start()

	jobChannel := <-pool

	t.Log("Creating a single job and send it to the worker")
	job := Job{GetTestEvent(3423543), 1}
	exp, _ = dialects.ConvertJSON(job.Event)
	jobChannel <- &job
	jobChannel = <-pool

	t.Log("Creating an another single job and send it to the worker")
	job = Job{GetTestEvent(1321), 1}
	exp, _ = dialects.ConvertJSON(job.Event)
	jobChannel <- &job
	jobChannel = <-pool

	t.Log("Send something that will fail and raise an error")
	sResp = fmt.Errorf("Error was intialized for testing")
	job = Job{GetTestEvent(43233), 1}
	exp, _ = dialects.ConvertJSON(job.Event)
	jobChannel <- &job
	jobChannel = <-pool

	if job.Attempt != 2 {
		t.Errorf("Job attempt number should be %d and it was %d instead", 2, job.Attempt)
	}

	t.Log("This failed message must be in the jobQueue, try again.")
	if len(jobQueue) != 1 {
		t.Errorf("jobChannel doesn't contain the previous job")
	}
	jobq := <-jobQueue
	sResp = nil
	jobChannel <- jobq
	jobChannel = <-pool

	t.Log("Send something that will fail and raise an error again")
	sResp = fmt.Errorf("Error was intialized for testing")
	job = Job{GetTestEvent(43254534), 1}
	exp, _ = dialects.ConvertJSON(job.Event)
	jobChannel <- &job
	jobChannel = <-pool

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

	var wg sync.WaitGroup
	worker.Stop(&wg)
}

// Worker Id testing
func TestGetId(t *testing.T) {
	jobQueue = make(chan *Job, 10)
	storageClient = &SimpleStorageClient{}

	t.Log("Creating a worker with 312 id")
	pool := make(chan chan *Job, 1)
	worker := NewWorker(312, 10, pool)

	if worker.GetId() != 312 {
		t.Errorf("Expected worker's ID was %d but it was %d instead.", 312, worker.GetId())
	}
}
