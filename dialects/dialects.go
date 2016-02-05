package dialects

import (
	"bytes"
)

// Interface for processing events
type StorageClient interface {
	IsBufferedStorage() bool
	GetConverter() Converter
	GetBatchConverter() BatchConverter
	Save(*bytes.Buffer) error
}

// Dialect interface for create StorageQueue from Configuration
type Dialect interface {
	IsValid() bool
	NewClient() (StorageClient, error)
}
