package main

import (
	"reflect"
	"testing"
)

// Testing the GetAction function's behaviour
func TestFunctionGetAction(t *testing.T) {
	t.Log("Testing the action type")
	// Testing event action
	eventAction := &EventAction{GetTestEvent(3423897841), 1}
	if exp := 1; eventAction.GetAction() != exp {
		t.Errorf("Expected action is %s but it was %s instead", exp, eventAction.GetAction())
	}
	// Testing flush action
	flushAction := &FlushAction{1}
	if exp := 2; flushAction.GetAction() != exp {
		t.Errorf("Expected masked IP setting is %s but it was %s instead", exp, flushAction.GetAction())
	}
}

// Testing the IsTargeted function's behaviour
func TestFunctionIsTargeted(t *testing.T) {
	t.Log("Testing action is targeted")
	// Testing event action
	eventAction := &EventAction{GetTestEvent(3423897841), 1}
	if exp := false; eventAction.IsTargeted() != exp {
		t.Errorf("Expected targeted is %s but it was %s instead", exp, eventAction.IsTargeted())
	}
	// Testing flush action
	flushAction := &FlushAction{3}
	if exp := true; flushAction.IsTargeted() != exp {
		t.Errorf("Expected targeted is %s but it was %s instead", exp, flushAction.IsTargeted())
	}
}

// Testing the GetTargetWorkerID function's behaviour
func TestFunctionGetTargetWorkerID(t *testing.T) {
	t.Log("Testing the target worker id")
	// Testing event action
	eventAction := &EventAction{GetTestEvent(3423897841), 1}
	if exp := -1; eventAction.GetTargetWorkerID() != exp {
		t.Errorf("Expected target worker id is %s but it was %s instead", exp, eventAction.GetTargetWorkerID())
	}
	// Testing flush action
	flushAction := &FlushAction{5}
	if exp := 5; flushAction.GetTargetWorkerID() != exp {
		t.Errorf("Expected target worker id is %s but it was %s instead", exp, flushAction.GetTargetWorkerID())
	}
}

// Testing the GetEvent function's behaviour
func TestFunctionGetEvent(t *testing.T) {
	t.Log("Testing action get event")
	// Testing event action
	event := GetTestEvent(3423897841)
	eventAction := &EventAction{event, 1}
	if !reflect.DeepEqual(GetTestEvent(3423897841), eventAction.GetEvent()) {
		t.Error("Not expected event was returned")
	}
}

// Testing the MarkAsFailed function's behaviour
func TestFunctionMarkAsFailed(t *testing.T) {
	t.Log("Testing mark as failed behaviour for jobs")
	jobQueue = make(chan Job, 10)
	job := &EventAction{GetTestEvent(3423897841), 1}

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
