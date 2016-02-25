package main

import (
	"testing"
)

// Testing the MarkAsFailed function's behaviour
func TestFunctionMarkAsFailed(t *testing.T) {
	t.Log("Testing mark as failed behaviour for jobs")
	jobQueue = make(chan *Job, 10)
	job := &Job{GetTestEvent(3423897841), 1}

	cases := []struct {
		ExpectedAttempt        int
		ExpectedJobQueueLength int
	}{
		{2, 1},
		{3, 2},
		{4, 2}}

	for i, c := range cases {
		t.Logf("Evaluate %d. attempt", i+1)
		job.MarkAsFailed(3)
		if job.Attempt != c.ExpectedAttempt {
			t.Errorf("Expected job's attempt was %d but it was %d instead", c.ExpectedAttempt, job.Attempt)
		}
		if len(jobQueue) != c.ExpectedJobQueueLength {
			t.Errorf("Expected jobQueue length was %d but it was %d instead", c.ExpectedJobQueueLength, len(jobQueue))
		}
	}
}
