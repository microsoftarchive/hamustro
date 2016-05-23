package main

import (
	"flag"
	"fmt"
	"github.com/wunderlist/hamustro/src/dialects"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var config *Config
var jobQueue chan Job
var storageClient dialects.StorageClient
var verbose bool
var isTerminating = false
var signatureRequired bool
var dispatcher *Dispatcher
var Version string = "1.0dev" // Current version

// Runs before the program starts
func main() {
	// Parse the CLI's attributes
	var filename = flag.String("config", "", "configuration `file` for the dialect")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode for debugging")
	flag.Parse()

	if *filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set a prefix for the logger
	log.SetPrefix(fmt.Sprintf("hamustro-%s ", Version))

	// Read and parse the configuration file
	config = NewConfig(*filename)
	if !config.IsValid() {
		log.Fatalf("Config is incomplete, please define `dialect` and `shared_secret` property")
	}

	// Set the signatureRequired variable
	signatureRequired = config.IsSignatureRequired()

	dialect, err := config.DialectConfig()
	if err != nil {
		log.Fatalf("Loading dialect configuration is failed: %s", err.Error())
	}
	if !dialect.IsValid() {
		log.Fatalf("Dialect configuration is incorrect or incomplete")
	}

	// Construct the dialect's client
	storageClient, err = dialect.NewClient()
	if err != nil {
		log.Fatalf("Client initialization is failed: %s", err.Error())
	}

	// Creates a worker options
	options := &WorkerOptions{
		BufferSize:   config.GetBufferSize(),
		RetryAttempt: config.GetRetryAttempt(),
		SpreadBuffer: config.IsSpreadBuffer()}

	// Create the background workers
	jobQueue = make(chan Job, config.GetMaxQueueSize())
	dispatcher = NewDispatcher(config.GetMaxWorkerSize(), options)
	dispatcher.Run()

	// Capture SIGINT and SIGTERM events to finish the ongoing work
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)
	go func() {
		<-signalChannel
		cleanup()
		os.Exit(1)
	}()

	// Set the log's output
	if config.LogFile != "" {
		logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Can't open logfile %s", err.Error())
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Start the server
	log.Printf("Starting server at %s", config.GetAddress())
	http.HandleFunc("/api/v1/track", TrackHandler)
	http.HandleFunc("/api/health", HealthHandler)
	http.HandleFunc("/api/flush", FlushHandler)
	if err := http.ListenAndServe(config.GetAddress(), nil); err != nil {
		log.Fatal(err)
	}
}

// Runs after the server was shut down
func cleanup() {
	// Do not accept new requests
	isTerminating = true
	log.Println("Shutting down server ...")

	// Set a timeout interval to force stop (avoid hanging out)
	go func() {
		c := time.Tick(90 * time.Second)
		for range c {
			log.Fatal("Server shut down is taking too long, force quit immediately.")
		}
	}()

	// Try to stop every worker
	dispatcher.Stop()
}
