package dialects

import (
	"encoding/json"
	"github.com/wunderlist/hamustro/src/payload"
	"regexp"
	"strconv"
)

func ConvertToJson(paramters []*payload.Parameter) string {
	out := map[string]string{}
	for _, p := range paramters {
		out[*p.Name] = *p.Value
	}
	b, _ := json.Marshal(out)
	return string(b[:])
}

// Single event
type Event struct {
	DeviceID        string `json:"device_id"`
	ClientID        string `json:"client_id"`
	Session         string `json:"session"`
	Nr              uint32 `json:"nr"`
	Env             string `json:"env"`
	SystemVersion   string `json:"system_version"`
	ProductVersion  string `json:"product_version"`
	At              string `json:"at"`
	Timezone        string `json:"timezone"`
	Event           string `json:"event"`
	DeviceMake      string `json:"device_make,omitempty"`
	DeviceModel     string `json:"device_model,omitempty"`
	System          string `json:"system,omitempty"`
	SystemLanguage  string `json:"system_language,omitempty"`
	Browser         string `json:"browser,omitempty"`
	BrowserVersion  string `json:"browser_version,omitempty"`
	ProductGitHash  string `json:"product_git_hash,omitempty"`
	ProductLanguage string `json:"product_language,omitempty"`
	UserID          string `json:"user_id,omitempty"`
	TenantID        string `json:"tenant_id,omitempty"`
	IP              string `json:"ip,omitempty"`
	Country         string `json:"country,omitempty"`
	Parameters      string `json:"parameters,omitempty"`
}

// Creates a new event based on the collection and a single payload
func NewEvent(meta *payload.Collection, payload *payload.Payload) *Event {
	return &Event{
		DeviceID:        meta.GetDeviceId(),
		ClientID:        meta.GetClientId(),
		Session:         meta.GetSession(),
		Nr:              payload.GetNr(),
		Env:             meta.GetEnv().String(),
		SystemVersion:   meta.GetSystemVersion(),
		ProductVersion:  meta.GetProductVersion(),
		At:              ConvertIsoformat(payload.GetAt()),
		Timezone:        payload.GetTimezone(),
		Event:           payload.GetEvent(),
		DeviceMake:      meta.GetDeviceMake(),
		DeviceModel:     meta.GetDeviceModel(),
		System:          meta.GetSystem(),
		SystemLanguage:  meta.GetSystemLanguage(),
		Browser:         meta.GetBrowser(),
		BrowserVersion:  meta.GetBrowserVersion(),
		ProductGitHash:  meta.GetProductGitHash(),
		ProductLanguage: meta.GetProductLanguage(),
		UserID:          payload.GetUserId(),
		TenantID:        payload.GetTenantId(),
		IP:              payload.GetIp(),
		Country:         payload.GetCountry(),
		Parameters:      ConvertToJson(payload.GetParameters())}
}

// Set IP Address
func (event *Event) SetIPAddress(IP string) {
	event.IP = IP
}

// Truncates the IP address
func (event *Event) TruncateIPv4LastOctet() {
	event.IP = regexpIP.ReplaceAllString(event.IP, "$1.0")
}

// Returns a
func (event *Event) String() []string {
	return []string{
		event.DeviceID,
		event.ClientID,
		event.Session,
		strconv.FormatInt(int64(event.Nr), 10),
		event.Env,
		event.SystemVersion,
		event.ProductVersion,
		event.At,
		event.Timezone,
		event.Event,
		event.DeviceMake,
		event.DeviceModel,
		event.System,
		event.SystemLanguage,
		event.Browser,
		event.BrowserVersion,
		event.ProductGitHash,
		event.ProductLanguage,
		event.UserID,
		event.TenantID,
		event.IP,
		event.Country,
		event.Parameters}
}

var regexpIP *regexp.Regexp

func init() {
	regexpIP, _ = regexp.Compile("([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})\\..+")
}
