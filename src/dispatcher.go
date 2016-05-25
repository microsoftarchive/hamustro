package main

import (
	"log"
	"sync"
	"time"
)

// A pool of workers channels that are registered with the dispatcher.
type Dispatcher struct {
	WorkerPool    chan *Worker
	Workers       []*Worker
	MaxWorkers    int
	WorkerOptions *WorkerOptions
}

// Options for worker creation
type FlushOptions struct {
	Automatic bool
}

// Creates a new dispatcher to handle new job requests
func NewDispatcher(maxWorkers int, options *WorkerOptions) *Dispatcher {
	pool := make(chan *Worker, maxWorkers)
	return &Dispatcher{
		WorkerPool:    pool,
		WorkerOptions: options,
		MaxWorkers:    maxWorkers}
}

// Returns the buffer size for a single worker
func (d *Dispatcher) GetBufferSize(n int) int {
	if !d.WorkerOptions.SpreadBuffer {
		return d.WorkerOptions.BufferSize
	}
	slizeSize := int(d.WorkerOptions.BufferSize / (2 * (d.MaxWorkers - 1)))
	return int(float32(d.WorkerOptions.BufferSize)*0.75) + (n * slizeSize)
}

// Creates and starts the workers
func (d *Dispatcher) Start() {
	for i := 0; i < d.MaxWorkers; i++ {
		options := &WorkerOptions{
			BufferSize:   d.GetBufferSize(i),
			RetryAttempt: d.WorkerOptions.RetryAttempt}

		// Create a new worker
		worker := NewWorker(i, options, d.WorkerPool)
		worker.Start()

		// Add the worker into the list
		d.Workers = append(d.Workers, worker)
	}
}

// Start automatic flush process
func (d *Dispatcher) TickAutomaticFlush() {
	tickerInterval := 60
	if config.AutoFlushInterval < 60 {
		tickerInterval = config.AutoFlushInterval
	}
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Flush(&FlushOptions{Automatic: true})
			}
		}
	}()
}

// Flush all the workers
func (d *Dispatcher) Flush(o *FlushOptions) {
	for i := range d.Workers {
		if o.Automatic == true && time.Now().Before(d.Workers[i].GetNextAutomaticFlush()) {
			continue
		}
		jobQueue <- &FlushAction{d.Workers[i].ID}
	}
}

// Creates and starts the workers and listen for new job requests
func (d *Dispatcher) Run() {
	d.Start()
	if config.AutoFlushInterval != 0 {
		d.TickAutomaticFlush()
	}
	go d.dispatch()
}

// Stops all the workers
func (d *Dispatcher) Stop() {
	var wg sync.WaitGroup
	for i := range d.Workers {
		wg.Add(1)
		d.Workers[i].Stop(&wg)
	}
	wg.Wait()
}

// Send the selected job to the worker
func (d *Dispatcher) Send(job Job, attempt int64) {
	select {
	case worker := <-d.WorkerPool:
		if !job.IsTargeted() || worker.ID == job.GetTargetWorkerID() {
			worker.JobChannel <- job
			return
		} else {
			if verbose {
				log.Printf("Re-register worker because targeted job arrived (for %d worker) but got %d worker instead", job.GetTargetWorkerID(), worker.ID)
			}
			go func() {
				// Try it a bit later after every failure.
				time.Sleep(time.Duration(100*attempt) * time.Millisecond)
				d.Send(job, attempt+1)
			}()
			d.WorkerPool <- worker
		}
	}
}

// Listening for new job requests
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-jobQueue:
			d.Send(job, 0)
		}
	}
}
