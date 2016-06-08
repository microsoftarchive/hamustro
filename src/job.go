package main

import (
	"github.com/wunderlist/hamustro/src/dialects"
)

// Define the known job's action types
const ACTION_EVENT = 1
const ACTION_FLUSH = 2

// Define the interface for Jobs
type Job interface {
	GetAction() int
	IsTargeted() bool
	GetTargetWorkerID() int
}

// Job: Flush action
type FlushAction struct {
	TargetWorkerID int
}

// Returns the name of the flush action
func (a *FlushAction) GetAction() int {
	return ACTION_FLUSH
}

// Returns that it's a targeted job
func (a *FlushAction) IsTargeted() bool {
	return true
}

// Returns the selected worker to do the action
func (a *FlushAction) GetTargetWorkerID() int {
	return a.TargetWorkerID
}

// Job: Add new Event action
type EventAction struct {
	Event   *dialects.Event
	Attempt int
}

// Returns the name of the add event action
func (a *EventAction) GetAction() int {
	return ACTION_EVENT
}

// Returns that it's not a targeted job
func (a *EventAction) IsTargeted() bool {
	return false
}

// Returns 0 to choose the first available worker
func (a *EventAction) GetTargetWorkerID() int {
	return -1
}

// Returns the Event object we'd like to add
func (a *EventAction) GetEvent() *dialects.Event {
	return a.Event
}

// Mark this job as failed and put back into the queue
func (a *EventAction) MarkAsFailed(retryAttempt int) {
	a.Attempt++
	if a.Attempt <= retryAttempt {
		jobQueue <- a
	}
}
