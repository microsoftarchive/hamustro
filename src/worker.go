package main

import (
	"fmt"
	"github.com/wunderlist/hamustro/src/dialects"
	"log"
	"sync"
	"time"
)

// Worker that executes the job.
type Worker struct {
	ID             int
	WorkerPool     chan *Worker
	JobChannel     chan Job
	BufferSize     int
	BufferedEvents []*dialects.Event
	Penalty        float32
	RetryAttempt   int
	LastSave       time.Time
	quit           chan *sync.WaitGroup
}

// Options for worker creation
type WorkerOptions struct {
	BufferSize   int
	RetryAttempt int
	SpreadBuffer bool
}

// Creates a new worker
func NewWorker(id int, options *WorkerOptions, workerPool chan *Worker) *Worker {
	return &Worker{
		ID:             id,
		WorkerPool:     workerPool,
		JobChannel:     make(chan Job),
		BufferSize:     options.BufferSize,
		BufferedEvents: []*dialects.Event{},
		Penalty:        1.0,
		RetryAttempt:   options.RetryAttempt,
		LastSave:       time.Now(),
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
			w.WorkerPool <- w

			select {
			case action := <-w.JobChannel:
				switch action.GetAction() {
				case ACTION_EVENT:
					if verbose {
						log.Printf("(%d worker) Received an add new event request!\n", w.ID)
					}
					if err := w.Work(action.(*EventAction)); err != nil {
						log.Print(err)
					}
				case ACTION_FLUSH:
					if verbose {
						log.Printf("(%d worker) Received a flush request!\n", w.ID)
					}
					if err := w.Flush(); err != nil {
						log.Print(err)
					}
				}
			case wg := <-w.quit:
				defer wg.Done()
				if err := w.Rescue(); err != nil {
					log.Print(err)
				}
				return
			}
		}
	}()
}

// Work on a single job
func (w *Worker) Work(action *EventAction) error {
	if !storageClient.IsBufferedStorage() {
		// Save messages
		if err := w.Save(action); err != nil {
			return err
		}

	} else {
		// Add message to the buffer if the storge is a buffered writer
		w.AddEventToBuffer(action.GetEvent())

		// Continue if the buffer is not full
		if !w.IsBufferFull() {
			return nil
		}

		if verbose {
			log.Printf("(%d worker) Saving buffered messages started", w.ID)
		}

		// Save messages
		if err := w.SaveBatch(); err != nil {
			return err
		}

		if verbose {
			log.Printf("(%d worker) Saving buffered messages was finished", w.ID)
		}
	}
	return nil
}

// Save Buffered messages
func (w *Worker) SaveBatch() error {
	// Convert messages to stringdefer
	msg, err := storageClient.GetBatchConverter()(w.BufferedEvents)
	if err != nil {
		w.IncreasePenalty()
		return fmt.Errorf("(%d worker) Batch converting buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
	}
	// Save messages
	if err := storageClient.Save(msg); err != nil {
		w.IncreasePenalty()
		return fmt.Errorf("(%d worker) Saving buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
	}
	w.ResetBuffer()
	w.UpdateLastSave()
	return nil
}

// Save messages
func (w *Worker) Save(action *EventAction) error {
	// Convert messages to string
	msg, err := storageClient.GetConverter()(action.GetEvent())
	if err != nil {
		rerr := fmt.Errorf("(%d worker) Encoding message is failed (%d attempt): %s", w.ID, action.Attempt, err.Error())
		action.MarkAsFailed(w.RetryAttempt)
		return rerr
	}

	// Save message immediately.
	if err := storageClient.Save(msg); err != nil {
		rerr := fmt.Errorf("(%d worker) Saving message is failed (%d attempt): %s", w.ID, action.Attempt, err.Error())
		action.MarkAsFailed(w.RetryAttempt)
		return rerr
	}
	w.UpdateLastSave()
	return nil
}

// Before the worker will be stopped it tries to rescue all ongoing job
func (w *Worker) Rescue() error {
	log.Printf("(%d worker) Received a signal to stop", w.ID)

	if err := w.Flush(); err != nil {
		return err
	}

	log.Printf("(%d worker) Stopped successfully", w.ID)
	return nil
}

// Flushing a worker
func (w *Worker) Flush() error {
	// We have received a signal to stop.
	if storageClient.IsBufferedStorage() && len(w.BufferedEvents) != 0 {
		log.Printf("(%d worker) Flushing %d buffered messages", w.ID, len(w.BufferedEvents))

		// Save messages
		if err := w.SaveBatch(); err != nil {
			return err
		}
	} else if len(w.BufferedEvents) == 0 {
		w.UpdateLastSave()
	}
	return nil
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

// Update last save time
func (w *Worker) UpdateLastSave() {
	w.LastSave = time.Now()
}

// Returns next possible automatic flush time
func (w *Worker) GetNextAutomaticFlush() time.Time {
	return w.LastSave.Add(time.Duration(config.AutoFlushInterval) * time.Second)
}

// Returns the worker's ID
func (w *Worker) GetId() int {
	return w.ID
}
