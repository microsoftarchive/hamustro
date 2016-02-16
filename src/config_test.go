package main

import (
	"os"
	"runtime"
	"testing"
)

// Testing a configuration loading from a file
func TestNewConfig(t *testing.T) {
}

// Testing worker size calculation
func TestGetMaxWorkerSize(t *testing.T) {
	t.Log("Testing worker size initialization")
	config := &Config{MaxWorkerSize: 0}
	if r := config.GetMaxWorkerSize(); r != runtime.NumCPU()+1 {
		t.Errorf("Expected worker size from default value was %d but it was %d instead", runtime.NumCPU()+1, r)
	}
	config = &Config{MaxWorkerSize: 433}
	if r := config.GetMaxWorkerSize(); r != 433 {
		t.Errorf("Expected worker size from configuration was %d but it was %d instead", 433, r)
	}
	os.Setenv("HAMUSTRO_MAX_WORKER_SIZE", "22")
	defer os.Unsetenv("HAMUSTRO_MAX_WORKER_SIZE")
	if r := config.GetMaxWorkerSize(); r != 22 {
		t.Errorf("Expected worker size from environment variable was %d but it was %d instead", 22, r)
	}
}

// Testing queue size calculation
func TestGetMaxQueueSize(t *testing.T) {
	t.Log("Testing queue size initialization")
	config := &Config{MaxQueueSize: 0}
	if r := config.GetMaxQueueSize(); r != (runtime.NumCPU()+1)*20 {
		t.Errorf("Expected queue size from default value was %d but it was %d instead", (runtime.NumCPU()+1)*20, r)
	}
	config = &Config{MaxQueueSize: 433}
	if r := config.GetMaxQueueSize(); r != 433 {
		t.Errorf("Expected queue size from configuration was %d but it was %d instead", 433, r)
	}
	os.Setenv("HAMUSTRO_MAX_QUEUE_SIZE", "22")
	defer os.Unsetenv("HAMUSTRO_MAX_QUEUE_SIZE")
	if r := config.GetMaxQueueSize(); r != 22 {
		t.Errorf("Expected queue size from environment variable was %d but it was %d instead", 22, r)
	}
}

// Testing port determination
func TestGetPort(t *testing.T) {
	t.Log("Testing port initialization")
	config := &Config{}
	if r := config.GetPort(); r != "8080" {
		t.Errorf("Expected port was %s but it was %s instead", "8080", r)
	}
	os.Setenv("HAMUSTRO_PORT", "8000")
	defer os.Unsetenv("HAMUSTRO_PORT")
	if r := config.GetPort(); r != "8000" {
		t.Errorf("Expected port was %s but it was %s instead", "8000", r)
	}
}

// Testing host determination
func TestGetHost(t *testing.T) {
	t.Log("Testing host initialization")
	config := &Config{}
	if r := config.GetHost(); r != "localhost" {
		t.Errorf("Expected host was %s but it was %s instead", "localhost", r)
	}
	os.Setenv("HAMUSTRO_HOST", "127.0.0.1")
	defer os.Unsetenv("HAMUSTRO_HOST")
	if r := config.GetHost(); r != "127.0.0.1" {
		t.Errorf("Expected host was %s but it was %s instead", "127.0.0.1", r)
	}
}

// Testing address determination
func TestGetAddress(t *testing.T) {
	t.Log("Testing address initialization")
	config := &Config{}
	if r := config.GetAddress(); r != "localhost:8080" {
		t.Errorf("Expected address was %s but it was %s instead", "localhost:8080", r)
	}

	os.Setenv("HAMUSTRO_PORT", "8000")
	defer os.Unsetenv("HAMUSTRO_PORT")
	if r := config.GetAddress(); r != "localhost:8000" {
		t.Errorf("Expected address was %s but it was %s instead", "localhost:8000", r)
	}

	os.Setenv("HAMUSTRO_HOST", "127.0.0.1")
	defer os.Unsetenv("HAMUSTRO_HOST")
	if r := config.GetAddress(); r != "127.0.0.1:8000" {
		t.Errorf("Expected address was %s but it was %s instead", "127.0.0.1:8000", r)
	}
}

// Testing the dialect determination
func TestDialectConfig(t *testing.T) {

}
