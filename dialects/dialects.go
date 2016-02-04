package dialects

// Interface for a worker that using a BufferedStorage as Client
type BufferedWorker interface {
	IsBufferFull() bool
	JoinBufferedEvents() *string
	ResetBuffer()
	AddEventToBuffer(*Event)
	GetId() int
}

// Interface for processing events
type StorageClient interface {
	IsBufferedStorage() bool
	Save(*string) error
}

// Dialect interface for create StorageQueue from Configuration
type Dialect interface {
	IsValid() bool
	NewClient() (StorageClient, error)
}
