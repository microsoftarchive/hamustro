package main

import (
	"hamustro/dialects"
)

// Job represents the job to be run.
type Job struct {
	Event   *dialects.Event
	Attempt int
}

// Mark this job as failed and put back into the queue
func (job *Job) MarkAsFailed(retryAttempt int) {
	job.Attempt++
	if job.Attempt <= retryAttempt {
		jobQueue <- job
	}
}
