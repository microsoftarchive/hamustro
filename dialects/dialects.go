package dialects

// Interface for processing events
type StorageClient interface {
	IsBufferedStorage() bool
	GetConverter() Converter
	Save(*string) error
}

// Dialect interface for create StorageQueue from Configuration
type Dialect interface {
	IsValid() bool
	NewClient() (StorageClient, error)
}
