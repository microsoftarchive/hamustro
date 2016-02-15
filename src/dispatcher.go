package main

import (
	"sync"
)

// A pool of workers channels that are registered with the dispatcher.
type Dispatcher struct {
	WorkerPool       chan chan *Job
	Workers          []*Worker
	MaxWorkers       int
	BufferSize       int
	SpreadBufferSize bool
}

func NewDispatcher(maxWorkers int, bufferSize int, spreadBufferSize bool) *Dispatcher {
	pool := make(chan chan *Job, maxWorkers)
	return &Dispatcher{
		WorkerPool:       pool,
		MaxWorkers:       maxWorkers,
		BufferSize:       bufferSize,
		SpreadBufferSize: spreadBufferSize}
}

func (d *Dispatcher) GetBufferSize(n int) int {
	if !d.SpreadBufferSize {
		return d.BufferSize
	}
	slizeSize := int(d.BufferSize / (2 * (d.MaxWorkers - 1)))
	return int(float32(d.BufferSize)*0.75) + (n * slizeSize)
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i, d.GetBufferSize(i), d.WorkerPool)
		worker.Start()
		d.Workers = append(d.Workers, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) Stop() {
	var wg sync.WaitGroup
	for i := range d.Workers {
		wg.Add(1)
		d.Workers[i].Stop(&wg)
	}
	wg.Wait()
}

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
