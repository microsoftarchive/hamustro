package main

import (
	"./dialects"
	"./dialects/abs"
	"./dialects/aqs"
	"./dialects/s3"
	"./dialects/sns"
	"./payload"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
)

// Job represents the job to be run.
type Job struct {
	Event   *dialects.Event
	Attempt int
}

// Mark this job as failed and put back into the queue
func (job *Job) MarkAsFailed() {
	job.Attempt++
	if job.Attempt <= 3 {
		jobQueue <- job
	}
}

// A buffered channel that we can send requests.
var JobQueue chan Job

// Worker that executes the job.
type Worker struct {
	ID             int
	WorkerPool     chan chan *Job
	JobChannel     chan *Job
	BufferSize     int
	BufferedEvents []*dialects.Event
	Penalty        float32
	quit           chan bool
}

func NewWorker(id int, bufferSize int, workerPool chan chan *Job) *Worker {
	return &Worker{
		ID:             id,
		WorkerPool:     workerPool,
		JobChannel:     make(chan *Job),
		BufferSize:     bufferSize,
		BufferedEvents: []*dialects.Event{},
		Penalty:        1.0,
		quit:           make(chan bool)}
}

// Start method starts the run loop for the worker.
// Listening for a quit channel in case we need to stop it.
func (w *Worker) Start() {
	go func() {
		for {
			// Register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// We have received a work request.
				if verbose {
					fmt.Printf("[%d] Received a job request!\n", w.ID)
				}
				if !storageClient.IsBufferedStorage() {
					// Convert the message to JSON string
					// TODO: Every dialect can define an output format!
					msg, err := job.Event.GetJSONMessage()
					if err != nil {
						log.Printf("(%d worker) Encoding message to JSON is failed (%d attempt): %s", w.ID, job.Attempt, err.Error())
						job.MarkAsFailed()
						continue
					}

					// Save message immediately.
					if err := storageClient.Save(&msg); err != nil {
						log.Printf("(%d worker) Saving message is failed (%d attempt): %s", w.ID, job.Attempt, err.Error())
						job.MarkAsFailed()
					}
				} else {
					// Add message to the buffer if the storge is a buffered writer
					w.AddEventToBuffer(job.Event)
					if w.IsBufferFull() {
						if err := storageClient.Save(w.JoinBufferedEvents()); err != nil {
							w.IncreasePenalty()
							log.Printf("(%d worker) Saving buffered messages is failed with %d records: %s", w.ID, len(w.BufferedEvents), err.Error())
							continue
						}
						w.ResetBuffer()
					}
				}
			case <-w.quit:
				// We have received a signal to stop.
				// TODO: Save the buffer before exiting.
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// Joins the buffered messages
func (w *Worker) JoinBufferedEvents() *string {
	s := []string{}
	for _, event := range w.BufferedEvents {
		msg, _ := event.GetJSONMessage()
		s = append(s, msg)
	}
	concat := strings.Join(s, "")
	return &concat
}

// Increase the value of the penalty attribute
func (w *Worker) IncreasePenalty() {
	w.Penalty *= 1.5
}

// Checks the state of the buffer
func (w *Worker) IsBufferFull() bool {
	return len(w.BufferedEvents) >= int(float32(w.BufferSize)*w.Penalty)
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

// A pool of workers channels that are registered with the dispatcher.
type Dispatcher struct {
	WorkerPool chan chan *Job
	MaxWorkers int
	BufferSize int
}

func NewDispatcher(maxWorkers int, bufferSize int) *Dispatcher {
	pool := make(chan chan *Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		BufferSize: bufferSize}
}

func (d *Dispatcher) Run() {
	// Starting `n` number of workers.
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(i, d.BufferSize, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
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

// Return the request's signature
func GetSignature(body []byte, time string) string {
	bodyHash := md5.New()
	io.WriteString(bodyHash, string(body[:]))

	requestHash := md5.New()
	io.WriteString(requestHash, time)
	io.WriteString(requestHash, "|")
	io.WriteString(requestHash, hex.EncodeToString(bodyHash.Sum(nil)))
	io.WriteString(requestHash, "|")
	io.WriteString(requestHash, config.SharedSecret)
	return hex.EncodeToString(requestHash.Sum(nil))
}

// Controller for `/api/v1/track`
func TrackHandler(w http.ResponseWriter, r *http.Request) {
	// Ignore not POST messages.
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if verbose {
			log.Println("Request was dropped: the sending method was not POST")
		}
		return
	}

	// If the client did not send time, we ignore
	if r.Header.Get("X-Hamustro-Time") == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if verbose {
			log.Println("Request was dropped: `X-Hamustro-Time` header is missing")
		}
		return
	}

	// If the client did not send signature of the message, we ignore
	if r.Header.Get("X-Hamustro-Signature") == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if verbose {
			log.Println("Request was dropped: `X-Hamustro-Signature` header is missing")
		}
		return
	}

	// Read the requests body into a variable.
	body, _ := ioutil.ReadAll(r.Body)

	// Calculate the request's signature
	if r.Header.Get("X-Hamustro-Signature") != GetSignature(body, r.Header.Get("X-Hamustro-Time")) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if verbose {
			log.Println("Request was dropped: `X-Hamustro-Signature` is not valid")
		}
		return
	}

	// Read the body into protobuf decoding.
	collection := &payload.Collection{}
	if err := proto.Unmarshal(body, collection); err != nil {
		log.Printf("Protobuf unmarshal is failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Creates a Job and put into the JobQueue for processing.
	for _, payload := range collection.Payloads {
		job := Job{dialects.NewEvent(collection, payload), 1}
		jobQueue <- &job
	}

	// Returns with 200.
	w.WriteHeader(http.StatusOK)
}

var config *Config
var jobQueue chan *Job
var storageClient dialects.StorageClient

// Application configuration
type Config struct {
	Host          string     `json:"host"`
	Port          string     `json:"port"`
	LogFile       string     `json:"logfile"`
	Dialect       string     `json:"dialect"`
	MaxWorkerSize int        `json:"max_worker_size"`
	MaxQueueSize  int        `json:"max_queue_size"`
	BufferSize    int        `json:"buffer_size"`
	SharedSecret  string     `json:"shared_secret"`
	AQS           aqs.Config `json:"aqs"`
	SNS           sns.Config `json:"sns"`
	ABS           abs.Config `json:"abs"`
	S3            s3.Config  `json:"s3"`
}

// Creates a new configuration object
func NewConfig(filename string) *Config {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}

// Returns the maximum worker size
func (c *Config) GetMaxWorkerSize() int {
	size, _ := strconv.ParseInt(os.Getenv("HAMUSTRO_MAX_WORKER_SIZE"), 10, 0)
	if size != 0 {
		return int(size)
	}
	if c.MaxWorkerSize != 0 {
		return c.MaxWorkerSize
	}
	return 5
}

// Returns the maximum queue size
func (c *Config) GetMaxQueueSize() int {
	size, _ := strconv.ParseInt(os.Getenv("HAMUSTRO_MAX_QUEUE_SIZE"), 10, 0)
	if size != 0 {
		return int(size)
	}
	if c.MaxQueueSize != 0 {
		return c.MaxQueueSize
	}
	return 100
}

// Returns the port of the application
func (c *Config) GetPort() string {
	if port := os.Getenv("HAMUSTRO_PORT"); port != "" {
		return port
	}
	return "8080"
}

// Returns the host of the application
func (c *Config) GetHost() string {
	if port := os.Getenv("HAMUSTRO_HOST"); port != "" {
		return port
	}
	return "localhost"
}

// Returns the address of the application
func (c *Config) GetAddress() string {
	return c.GetHost() + ":" + c.GetPort()
}

func (c *Config) DialectConfig() (dialects.Dialect, error) {
	switch c.Dialect {
	case "aqs":
		return &c.AQS, nil
	case "sns":
		return &c.SNS, nil
	case "abs":
		return &c.ABS, nil
	case "s3":
		return &c.S3, nil
	}
	return nil, errors.New("not supported `dialect` in the configuration file.")
}

var verbose bool

func init() {
	var filename = flag.String("config", "", "configuration `file` for the dialect")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode for debugging")
	flag.Parse()

	if *filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config = NewConfig(*filename)
	dialect, err := config.DialectConfig()
	if err != nil {
		log.Fatalf("Loading dialect configuration is failed: %s", err.Error())
	}
	if !dialect.IsValid() {
		log.Fatalf("Dialect configuration is incorrect or incomplete: %s", err.Error())
	}
	storageClient, err = dialect.NewClient()
	if err != nil {
		log.Fatalf("Client initialization is failed: %s", err.Error())
	}

	jobQueue = make(chan *Job, config.GetMaxQueueSize())
	dispatcher := NewDispatcher(config.GetMaxWorkerSize(), config.BufferSize)
	dispatcher.Run()
}

func main() {
	if config.LogFile != "" {
		logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Can't open logfile %s", err.Error())
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	log.Printf("Starting server at %s", config.GetAddress())
	http.HandleFunc("/api/v1/track", TrackHandler)
	http.ListenAndServe(config.GetAddress(), nil)
}
