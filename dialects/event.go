package dialects

import (
	"../payload"
	"strconv"
	"time"
)

// Converts EPOCH timestamp to isoformat string
func ConvertIsoformat(at uint64) string {
	return time.Unix(int64(at), 0).Format("2006-01-02T15:04:05")
}

// Single event
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

// Creates a new event based on the collection and a single payload
func NewEvent(meta *payload.Collection, payload *payload.Payload) *Event {
	return &Event{
		DeviceID:       meta.GetDeviceId(),
		ClientID:       meta.GetClientId(),
		Session:        meta.GetSession(),
		Nr:             payload.GetNr(),
		SystemVersion:  meta.GetSystemVersion(),
		ProductVersion: meta.GetProductVersion(),
		At:             ConvertIsoformat(payload.GetAt()),
		Event:          payload.GetEvent(),
		System:         meta.GetSystem(),
		ProductGitHash: meta.GetProductGitHash(),
		UserID:         payload.GetUserId(),
		IP:             payload.GetIp(),
		Parameters:     payload.GetParameters(),
		IsTesting:      payload.GetIsTesting()}
}

// Returns a
func (event *Event) String() []string {
	return []string{
		event.DeviceID,
		event.ClientID,
		event.Session,
		strconv.FormatInt(int64(event.Nr), 10),
		event.SystemVersion,
		event.ProductVersion,
		event.At,
		event.Event,
		event.System,
		event.ProductGitHash,
		strconv.FormatInt(int64(event.UserID), 10),
		event.IP,
		event.Parameters,
		strconv.FormatBool(event.IsTesting)}
}
