package main

import (
	"encoding/json"
	"fmt"
	"github.com/wunderlist/hamustro/src/dialects"
	"github.com/wunderlist/hamustro/src/dialects/abs"
	"github.com/wunderlist/hamustro/src/dialects/aqs"
	"github.com/wunderlist/hamustro/src/dialects/s3"
	"github.com/wunderlist/hamustro/src/dialects/sns"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Application configuration
type Config struct {
	LogFile           string     `json:"logfile"`
	Dialect           string     `json:"dialect"`
	MaxWorkerSize     int        `json:"max_worker_size"`
	MaxQueueSize      int        `json:"max_queue_size"`
	RetryAttempt      int        `json:"retry_attempt"`
	BufferSize        int        `json:"buffer_size"`
	MaskedIP          bool       `json:"masked_ip"`
	SpreadBufferSize  bool       `json:"spread_buffer_size"`
	Signature         string     `json:"signature"`
	SharedSecret      string     `json:"shared_secret"`
	MaintenanceKey    string     `json:"maintenance_key"`
	AutoFlushInterval int        `json:"auto_flush_interval"`
	AQS               aqs.Config `json:"aqs"`
	SNS               sns.Config `json:"sns"`
	ABS               abs.Config `json:"abs"`
	S3                s3.Config  `json:"s3"`
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

// Configuration validation
func (c *Config) IsValid() bool {
	return c.Dialect != "" && c.SharedSecret != ""
}

// Get Signature's status
func (c *Config) IsSignatureRequired() bool {
	switch c.Signature {
	case "required":
		return true
	case "optional":
		return false
	default:
		return true
	}
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
	return runtime.NumCPU() + 1
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
	return c.GetMaxWorkerSize() * 20
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
	host := c.GetHost()
	if host == "localhost" {
		return ":" + c.GetPort()
	}
	return host + ":" + c.GetPort()
}

// Returns the default buffer size for Buffered Storage.
func (c *Config) GetBufferSize() int {
	if c.BufferSize != 0 {
		return c.BufferSize
	}
	return (c.GetMaxWorkerSize() * c.GetMaxQueueSize()) * 10
}

// Returns the default flush timeout
func (c *Config) GetAutoFlushInterval() int {
	if c.AutoFlushInterval != 0 {
		return c.AutoFlushInterval
	}
	return 0
}

// Returns the default spreding property
func (c *Config) IsSpreadBuffer() bool {
	return c.SpreadBufferSize
}

// Should we truncate the IP address
func (c *Config) IsMaskedIP() bool {
	return c.MaskedIP
}

// Returns the retry attempt number
func (c *Config) GetRetryAttempt() int {
	if c.RetryAttempt != 0 {
		return c.RetryAttempt
	}
	return 3
}

// Returns the selected dialect's configuration object
func (c *Config) DialectConfig() (dialects.Dialect, error) {
	switch strings.ToLower(c.Dialect) {
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
