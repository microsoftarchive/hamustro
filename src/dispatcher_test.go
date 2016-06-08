package main

import (
	"github.com/wunderlist/hamustro/src/dialects"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

// Testing retry attempt settings
func TestRetryAttemptNewDispatcher(t *testing.T) {
	t.Log("Testing retry attempt settings")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{RetryAttempt: 5}
	dispatcher := NewDispatcher(4, options)
	dispatcher.Start()
	if exp := 4; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}
	for _, w := range dispatcher.Workers {
		if exp := 5; w.RetryAttempt != exp {
			t.Errorf("Expected %d worker's retry attempt property was %d and it was %d instead", w.ID, exp, w.RetryAttempt)
		}
	}
}

// Testing buffer size calculation for not spreading buffers
func TestNotSpreadBufferNewDispatcher(t *testing.T) {
	t.Log("Testing buffer size calculation for not spreading buffers")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{SpreadBuffer: false, BufferSize: 10000}
	dispatcher := NewDispatcher(4, options)
	dispatcher.Start()
	if exp := 4; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}
	for _, w := range dispatcher.Workers {
		if exp := 10000; w.BufferSize != exp {
			t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", w.ID, exp, w.BufferSize)
		}
	}
}

// Testing buffer size calculation for not spreading buffers
func TestSpreadBufferNewDispatcher(t *testing.T) {
	t.Log("Testing buffer size calculation for not spreading buffers")
	storageClient = &SimpleStorageClient{}
	log.SetOutput(ioutil.Discard)

	options := &WorkerOptions{SpreadBuffer: true, BufferSize: 10000}
	dispatcher := NewDispatcher(3, options)
	dispatcher.Start()
	if exp := 3; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	cases := []int{7500, 10000, 12500}
	for i, exp := range cases {
		if dispatcher.Workers[i].BufferSize != exp {
			t.Errorf("Expected %d worker's buffer size was %d and it was %d instead", dispatcher.Workers[i].ID, exp, dispatcher.Workers[i].BufferSize)
		}
	}
}

// Testing the buffer size caluculation function
func TestFunctionDispatcherGetBufferSize(t *testing.T) {
	t.Log("Testing the buffer size caluculation function")
	dispatcher := &Dispatcher{MaxWorkers: 3, WorkerOptions: &WorkerOptions{BufferSize: 10000, SpreadBuffer: true}}
	for i, exp := range []int{7500, 10000, 12500} {
		if size := dispatcher.GetBufferSize(i); size != exp {
			t.Errorf("Expected buffer size was %d and it was %d instead", exp, size)
		}
	}

	dispatcher = &Dispatcher{MaxWorkers: 3, WorkerOptions: &WorkerOptions{BufferSize: 10000, SpreadBuffer: false}}
	for i, exp := range []int{10000, 10000, 10000} {
		if size := dispatcher.GetBufferSize(i); size != exp {
			t.Errorf("Expected buffer size was %d and it was %d instead", exp, size)
		}
	}
}

// Testing the dispatcher listen function
func TestDispatcherListen(t *testing.T) {
	config = &Config{}                     // Define an empty config
	storageClient = &SimpleStorageClient{} // Define the Simple Storage as a storage
	jobQueue = make(chan Job, 10)          //Define the job Queue
	log.SetOutput(ioutil.Discard)          // Disable the logger
	T, response, catched = t, nil, false   // Set properties

	t.Log("Creates the dispatcher and listen for new jobs")
	options := &WorkerOptions{RetryAttempt: 5}
	dispatcher := NewDispatcher(2, options)
	dispatcher.Run()

	if exp := 2; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	t.Log("Creating two jobs and put it into the job queue")
	job1 := EventAction{GetTestEvent(423432), 1}
	expBuffer1, _ := dialects.ConvertJSON(job1.Event)

	job2 := EventAction{GetTestEvent(7643329), 1}
	expBuffer2, _ := dialects.ConvertJSON(job2.Event)

	t.Log("It should catch a different worker with the expected results")
	exp = map[string]struct{}{expBuffer1.String(): {}, expBuffer2.String(): {}}
	jobQueue <- &job1
	jobQueue <- &job2

	t.Log("Wait until both is finished")
	time.Sleep(150 * time.Millisecond)
	if !catched {
		t.Errorf("Worker didn't catch the expected jobs")
	}

	t.Log("Stops the workers")
	dispatcher.Stop()
}

// Testing the dispatcher listen function
func TestDispatcherFlush(t *testing.T) {
	config = &Config{}                       // Define an empty config
	storageClient = &BufferedStorageClient{} // Define the Buffered Storage as a storage
	jobQueue = make(chan Job, 10)            //Define the job Queue
	log.SetOutput(ioutil.Discard)            // Disable the logger
	T, response, catched = t, nil, false     // Set properties

	t.Log("Creates the dispatcher with a single worker and listen for new jobs")
	options := &WorkerOptions{BufferSize: 3}
	dispatcher := NewDispatcher(1, options)
	dispatcher.Run()

	if exp := 1; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	t.Log("Creating a job and put it into the job queue")
	job := EventAction{GetTestEvent(636284), 1}
	expBuffer, _ := dialects.ConvertJSON(job.Event)

	exp = map[string]struct{}{expBuffer.String(): {}}
	jobQueue <- &job

	t.Log("No saving should happen because of the buffer size")
	time.Sleep(150 * time.Millisecond)
	if catched {
		t.Errorf("Worker shouldn't catch the job")
	}

	t.Log("Run an API flush")
	dispatcher.Flush(&FlushOptions{Automatic: false})

	t.Log("Wait until the flush is finished, saving should be emited and catched")
	time.Sleep(150 * time.Millisecond)
	if !catched {
		t.Errorf("Worker didn't catch the job")
	}
}

// Testing the dispatcher listen function
func TestDispatcherAutomaticFlush(t *testing.T) {
	config = &Config{AutoFlushInterval: 3}   // Define the config
	storageClient = &BufferedStorageClient{} // Define the Buffered Storage as a storage
	jobQueue = make(chan Job, 10)            //Define the job Queue
	log.SetOutput(ioutil.Discard)            // Disable the logger
	T, response, catched = t, nil, false     // Set properties

	t.Log("Creates the dispatcher with a single worker and listen for new jobs")
	options := &WorkerOptions{BufferSize: 3}
	dispatcher := NewDispatcher(1, options)
	dispatcher.Run()

	if exp := 1; len(dispatcher.Workers) != exp {
		t.Errorf("Expected worker's count was %d but it was %d instead", exp, len(dispatcher.Workers))
	}

	t.Log("Creating a job and put it into the job queue")
	job := EventAction{GetTestEvent(636284), 1}
	expBuffer, _ := dialects.ConvertJSON(job.Event)

	exp = map[string]struct{}{expBuffer.String(): {}}
	jobQueue <- &job

	t.Log("Wait until it's finished, it should be buffered and not saved")
	time.Sleep(150 * time.Millisecond)
	if catched {
		t.Errorf("Worker shouldn't catch the job")
	}

	t.Log("Run an automatic flush")
	dispatcher.Flush(&FlushOptions{Automatic: true})

	t.Log("Wait and check the flush isn't finished, it should not be fired")
	time.Sleep(150 * time.Millisecond)
	if catched {
		t.Errorf("Worker shouldn't catch the job")
	}

	t.Log("Wait until the interval is over (started by the dispatcher) and check again")
	time.Sleep(3 * time.Second)
	if !catched {
		t.Errorf("Worker didn't catch the job")
	}
}

// Testing the dispatcher listen function
func TestDispatcherWaitingForFlush(t *testing.T) {
	config = &Config{}                                      // Define an empty config
	storageClient = &BufferedStorageClientWithoutExpected{} // Define the Simple Storage as a storage
	jobQueue = make(chan Job, 10)                           // Define the job Queue
	log.SetOutput(ioutil.Discard)                           // Disable the logger
	T, response, catched = t, nil, false                    // Set properties

	pool := make(chan *Worker, 2)
	workerOptions := &WorkerOptions{BufferSize: 10}

	t.Log("Creates the dispatcher with two workers and listen for new jobs")
	dispatcher := &Dispatcher{
		WorkerPool:    pool,
		WorkerOptions: workerOptions,
		MaxWorkers:    2}

	t.Log("Create two workers and start them")
	worker1 := NewWorker(1, workerOptions, dispatcher.WorkerPool)
	worker1.Start()

	worker2 := NewWorker(2, workerOptions, dispatcher.WorkerPool)
	worker2.Start()

	t.Log("Append workers to the dispatcher and start the dispatcher")
	dispatcher.Workers = append(dispatcher.Workers, worker1)
	dispatcher.Workers = append(dispatcher.Workers, worker2)

	go dispatcher.dispatch()

	t.Log("Create two events, and send these to the workers")
	action1 := EventAction{GetTestEvent(33344), 1}
	action2 := EventAction{GetTestEvent(88829), 1}

	<-dispatcher.WorkerPool
	<-dispatcher.WorkerPool
	worker1.JobChannel <- &action1
	worker2.JobChannel <- &action2

	t.Log("Wait until the workers catch the jobs")
	time.Sleep(150 * time.Millisecond)

	CheckResultsForBufferedStorage(worker1, 1, 1.0, 10)
	CheckResultsForBufferedStorage(worker2, 1, 1.0, 10)

	t.Log("Stop one of the workers randomly")
	stopped_worker := <-dispatcher.WorkerPool
	running_worker := worker2
	if worker1.ID != stopped_worker.ID {
		running_worker = worker1
	}

	t.Log("Flush all of the workers")
	dispatcher.Flush(&FlushOptions{Automatic: false})

	t.Log("Wait until the flush finish for the active worker")
	time.Sleep(150 * time.Millisecond)

	ValidateSending()
	CheckResultsForBufferedStorage(running_worker, 0, 1.0, 10)
	CheckResultsForBufferedStorage(stopped_worker, 1, 1.0, 10)

	t.Log("Create two new event, and send these to the active worker")
	action3 := EventAction{GetTestEvent(11122), 1}
	action4 := EventAction{GetTestEvent(88765), 1}

	jobQueue <- &action3
	jobQueue <- &action4

	t.Log("Wait until the active worker is processing the jobs without saving")
	time.Sleep(200 * time.Millisecond)
	if catched {
		t.Errorf("Worker shouldn't catch the job")
	}

	CheckResultsForBufferedStorage(running_worker, 2, 1.0, 10)
	CheckResultsForBufferedStorage(stopped_worker, 1, 1.0, 10)

	t.Log("Start the previously stopped worker")
	dispatcher.WorkerPool <- stopped_worker

	t.Log("This worker should catch the forgotten flush job")
	time.Sleep(1 * time.Second)
	ValidateSending()
	CheckResultsForBufferedStorage(running_worker, 2, 1.0, 10)
	CheckResultsForBufferedStorage(stopped_worker, 0, 1.0, 10)

	t.Log("Stop the workers and flush the events from the active worker")
	var wg sync.WaitGroup
	wg.Add(2)
	worker1.Stop(&wg)
	worker2.Stop(&wg)
	wg.Wait()
	ValidateSending()
	CheckResultsForBufferedStorage(worker1, 0, 1.0, 10)
	CheckResultsForBufferedStorage(worker2, 0, 1.0, 10)
}
