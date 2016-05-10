package main

import (
	"sync"
	"time"
)

// A pool of workers channels that are registered with the dispatcher.
type Dispatcher struct {
	WorkerPool    chan chan *Job
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
	pool := make(chan chan *Job, maxWorkers)
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
func (d *Dispatcher) StartAutomaticFlush() {
	if config.GetAutoFlushInterval() == 0 {
		return
	}
	ticker := time.NewTicker(60 * time.Second)
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
		if d.Workers[i].IsSaving == true {
			continue
		}
		d.Workers[i].SetIsSaving(true)
		d.Workers[i].Flush(o)
		d.Workers[i].SetIsSaving(false)
	}
}

// Creates and starts the workers and listen for new job requests
func (d *Dispatcher) Run() {
	d.Start()
	d.StartAutomaticFlush()
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

// Listening for new job requests
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-jobQueue:
			select {
			case jobChannel := <-d.WorkerPool:
				jobChannel <- job
			}
		}
	}
}
