package main

import (
	"./dialects"
	"fmt"
	"log"
	"sync"
)

// Worker that executes the job.
type Worker struct {
	ID             int
	WorkerPool     chan chan *Job
	JobChannel     chan *Job
	BufferSize     int
	BufferedEvents []*dialects.Event
	Penalty        float32
	RetryAttempt   int
	quit           chan *sync.WaitGroup
}

// Options for worker creation
type WorkerOptions struct {
	BufferSize   int
	RetryAttempt int
	SpreadBuffer bool
}

// Creates a new worker
func NewWorker(id int, options *WorkerOptions, workerPool chan chan *Job) *Worker {
	return &Worker{
		ID:             id,
		WorkerPool:     workerPool,
		JobChannel:     make(chan *Job),
		BufferSize:     options.BufferSize,
		BufferedEvents: []*dialects.Event{},
		Penalty:        1.0,
		RetryAttempt:   options.RetryAttempt,
		quit:           make(chan *sync.WaitGroup)}
}

// Start method starts the run loop for the worker.
// Listening for a quit channel in case we need to stop it.
func (w *Worker) Start() {
	if storageClient.IsBufferedStorage() {
		log.Printf("(%d worker) Started with %d buffer", w.ID, w.BufferSize)
	} else {
		log.Printf("(%d worker) Started", w.ID)
	}
	go func() {
		for {
			// Register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// We have received a work request.
				if verbose {
					fmt.Printf("(%d worker) Received a job request!\n", w.ID)
				}
				if !storageClient.IsBufferedStorage() {
					// Convert the message to string
					msg, err := storageClient.GetConverter()(job.Event)
					if err != nil {
						log.Printf("(%d worker) Encoding message is failed (%d attempt): %s", w.ID, job.Attempt, err.Error())
						job.MarkAsFailed(w.RetryAttempt)
						continue
					}

					// Save message immediately.
					if err := storageClient.Save(w.ID, msg); err != nil {
						log.Printf("(%d worker) Saving message is failed (%d attempt): %s", w.ID, job.Attempt, err.Error())
						job.MarkAsFailed(w.RetryAttempt)
						continue
					}
				} else {
					// Add message to the buffer if the storge is a buffered writer
					w.AddEventToBuffer(job.Event)

					// Continue if the buffer is not full
					if !w.IsBufferFull() {
						continue
					}

					// Convert messages to string
					msg, err := storageClient.GetBatchConverter()(w.BufferedEvents)
					if err != nil {
						w.IncreasePenalty()
						log.Printf("(%d worker) Batch converting buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
						continue
					}
					// Save messages
					if err := storageClient.Save(w.ID, msg); err != nil {
						w.IncreasePenalty()
						log.Printf("(%d worker) Saving buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
						continue
					}
					w.ResetBuffer()
				}
			case wg := <-w.quit:
				defer wg.Done()
				log.Printf("(%d worker) Received a signal to stop", w.ID)

				// We have received a signal to stop.
				if storageClient.IsBufferedStorage() && len(w.BufferedEvents) != 0 {
					log.Printf("(%d worker) Flushing %d buffered messages", w.ID, len(w.BufferedEvents))

					// Convert messages to string
					msg, err := storageClient.GetBatchConverter()(w.BufferedEvents)
					if err != nil {
						log.Printf("(%d worker) Batch converting buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
						return
					}
					// Save messages
					if err := storageClient.Save(w.ID, msg); err != nil {
						log.Printf("(%d worker) Saving buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
						return
					}
					w.ResetBuffer()
				}
				log.Printf("(%d worker) Stopped successfully", w.ID)
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop(wg *sync.WaitGroup) {
	go func() {
		if storageClient.IsBufferedStorage() {
			log.Printf("(%d worker) Sending stop signal to worker with %d buffered events", w.ID, len(w.BufferedEvents))
		} else {
			log.Printf("(%d worker) Sending stop signal to worker", w.ID)
		}
		w.quit <- wg
	}()
}

// Increase the value of the penalty attribute
func (w *Worker) IncreasePenalty() {
	w.Penalty *= 1.5
}

// Returns the current buffer size with the current penalty
func (w *Worker) GetBufferSize() int {
	return int(float32(w.BufferSize) * w.Penalty)
}

// Checks the state of the buffer
func (w *Worker) IsBufferFull() bool {
	return len(w.BufferedEvents) >= w.GetBufferSize()
}

// Resets the buffer
func (w *Worker) ResetBuffer() {
	w.BufferedEvents = w.BufferedEvents[:0]
	w.Penalty = 1.0
}

// Adds a message to the buffer
func (w *Worker) AddEventToBuffer(event *dialects.Event) {
	w.BufferedEvents = append(w.BufferedEvents, event)
}

// Returns the worker's ID
func (w *Worker) GetId() int {
	return w.ID
}
