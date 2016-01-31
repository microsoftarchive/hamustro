package dialects

import (
	"../payload"
	"bytes"
	"encoding/json"
)

type Event struct {
	DeviceID       string `json:"device_id"`
	ClientID       string `json:"client_id"`
	Session        string `json:"session"`
	Nr             uint32 `json:"nr"`
	SystemVersion  string `json:"system_version"`
	ProductVersion string `json:"product_version"`
	At             string `json:"at"`
	Event          string `json:"event"`
	System         string `json:"system,omitempty"`
	ProductGitHash string `json:"product_git_hash,omitempty"`
	UserID         uint32 `json:"user_id,omitempty"`
	IP             string `json:"ip,omitempty"`
	Parameters     string `json:"parameters,omitempty"`
	IsTesting      bool   `json:"is_testing"`
}

func NewEvent(meta *payload.Collection, payload *payload.Payload) *Event {
	return &Event{
		DeviceID:       meta.GetDeviceId(),
		ClientID:       meta.GetClientId(),
		Session:        meta.GetSession(),
		Nr:             payload.GetNr(),
		SystemVersion:  meta.GetSystemVersion(),
		ProductVersion: meta.GetProductVersion(),
		At:             payload.GetAt(),
		Event:          payload.GetEvent(),
		System:         meta.GetSystem(),
		ProductGitHash: meta.GetProductGitHash(),
		UserID:         payload.GetUserId(),
		IP:             payload.GetIp(),
		Parameters:     payload.GetParameters(),
		IsTesting:      payload.GetIsTesting()}
}

// Send a single Event into the queue
func (event *Event) GetJSONMessage() (string, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(event); err != nil {
		return "", err
	}
	return b.String(), nil
}

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
