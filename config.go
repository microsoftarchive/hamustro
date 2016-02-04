package main

import (
	"./dialects"
	"./dialects/abs"
	"./dialects/aqs"
	"./dialects/s3"
	"./dialects/sns"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Application configuration
type Config struct {
	Host             string     `json:"host"`
	Port             string     `json:"port"`
	LogFile          string     `json:"logfile"`
	Dialect          string     `json:"dialect"`
	MaxWorkerSize    int        `json:"max_worker_size"`
	MaxQueueSize     int        `json:"max_queue_size"`
	BufferSize       int        `json:"buffer_size"`
	SpreadBufferSize bool       `json:"spread_buffer_size"`
	SharedSecret     string     `json:"shared_secret"`
	AQS              aqs.Config `json:"aqs"`
	SNS              sns.Config `json:"sns"`
	ABS              abs.Config `json:"abs"`
	S3               s3.Config  `json:"s3"`
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

// Returns the selected dialect's configuration object
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
	return nil, fmt.Errorf("Not supported `%s` dialect in the configuration file.", c.Dialect)
}
