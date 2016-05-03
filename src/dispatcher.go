package main

import (
	"sync"
	"log"
	"time"
)

// A pool of workers channels that are registered with the dispatcher.
type Dispatcher struct {
	WorkerPool    chan chan *Job
	Workers       []*Worker
	MaxWorkers    int
	WorkerOptions *WorkerOptions
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
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <- ticker.C:
				d.AutomaticFlush()
			}
		}
	}()
}

// Automatic flush all workers
func (d *Dispatcher) AutomaticFlush() {
	for i := range d.Workers {
		if err := d.Workers[i].AutomaticFlush(); err != nil {
			log.Print(err)
		}
	}
}

// Creates and starts the workers and listen for new job requests
func (d *Dispatcher) Run() {
	d.Start()
	d.StartAutomaticFlush()
	go d.dispatch()
}

// Flush all the workers
func (d *Dispatcher) Flush() {
	for i := range d.Workers {
		if err := d.Workers[i].Flush(); err != nil {
			log.Print(err)
		}
	}
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
