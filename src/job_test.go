package main

import (
	"testing"
)

// Testing the MarkAsFailed function's behaviour
func TestMarkAsFailed(t *testing.T) {
	t.Log("Testing mark as failed behaviour for jobs")
	jobQueue = make(chan *Job, 10)
	job := &Job{GetTestEvent(3423897841), 1}

	t.Log("Evaluate the first attempt")
	job.MarkAsFailed(3)
	if job.Attempt != 2 {
		t.Errorf("Expected job's attempt was %d but it was %d instead", 2, job.Attempt)
	}
	if len(jobQueue) != 1 {
		t.Errorf("Expected jobQueue length was %d but it was %d instead", 1, len(jobQueue))
	}

	t.Log("Evaluate the second attempt")
	jobQueue = make(chan *Job, 10)
	job.MarkAsFailed(3)
	if job.Attempt != 3 {
		t.Errorf("Expected job's attempt was %d but it was %d instead", 3, job.Attempt)
	}
	if len(jobQueue) != 1 {
		t.Errorf("Expected jobQueue length was %d but it was %d instead", 1, len(jobQueue))
	}

	t.Log("Evaluate the third attempt")
	jobQueue = make(chan *Job, 10)
	job.MarkAsFailed(3)
	if job.Attempt != 4 {
		t.Errorf("Expected job's attempt was %d but it was %d instead", 4, job.Attempt)
	}
	if len(jobQueue) != 0 {
		t.Errorf("Expected jobQueue length was %d but it was %d instead", 0, len(jobQueue))
	}
}
