package dialects

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/payload"
	"reflect"
	"testing"
)

// Returns an Event for testing purposes
func GetTestEvent(userId uint32) *Event {
	return &Event{
		DeviceID:        "a73b1c37-2c24-4786-af7a-16de88fbe23a",
		ClientID:        "bce44f67b2661fd445d469b525b04f68",
		Session:         "0e350a2cc31648bb24ba61eb14be337c",
		Nr:              1,
		Env:             "PRODUCTION",
		SystemVersion:   "10.10",
		ProductVersion:  "1.1.2",
		At:              "2016-02-05T15:05:04",
		Timezone:        "+02:00",
		Event:           "Client.CreateUser",
		DeviceMake:      "Iphone",
		DeviceModel:     "Iphone 6",
		System:          "OSX",
		SystemLanguage:  "DE",
		Browser:         "Mozilla",
		BrowserVersion:  "10.01.11",
		ProductGitHash:  "5416a5889392d509e3bafcf40f6388e83aab23e6",
		ProductLanguage: "HU",
		UserID:          fmt.Sprintf("%v", userId),
		TenantID:        "sdfghjkloiuytremiwoz",
		IP:              "214.160.227.22",
		Country:         "UK",
		Parameters:      "{\"parameter\": \"test_parameter\"}"}
}

// Converts and Event into a list of string
func TestEventStringConversion(t *testing.T) {
	t.Log("Converting Event into a list of strings")
	e := GetTestEvent(97421193)
	exp := []string{
		"a73b1c37-2c24-4786-af7a-16de88fbe23a",
		"bce44f67b2661fd445d469b525b04f68",
		"0e350a2cc31648bb24ba61eb14be337c",
		"1",
		"PRODUCTION",
		"10.10",
		"1.1.2",
		"2016-02-05T15:05:04",
		"+02:00",
		"Client.CreateUser",
		"Iphone",
		"Iphone 6",
		"OSX",
		"DE",
		"Mozilla",
		"10.01.11",
		"5416a5889392d509e3bafcf40f6388e83aab23e6",
		"HU",
		"97421193",
		"sdfghjkloiuytremiwoz",
		"214.160.227.22",
		"UK",
		"{\"parameter\": \"test_parameter\"}"}
	if !reflect.DeepEqual(e.String(), exp) {
		t.Error("Expected event's string is not matched")
	}
}

// Tests the IP truncating functionality
func TestFunctionTruncateIPv4LastOctet(t *testing.T) {
	cases := []struct {
		IP         string
		ExpectedIP string
	}{
		{"214.160.227.22", "214.160.227.0"},
		{"214.160.227.22:80", "214.160.227.0"},
		{"214.160.227.22/24", "214.160.227.0"},
		{"", ""}}
	t.Log("Removing last octet of IP address")
	for _, c := range cases {
		e := &Event{IP: c.IP}
		e.TruncateIPv4LastOctet()
		if c.ExpectedIP != e.IP {
			t.Errorf("Expected truncated IP was %s but it was %s instead", c.ExpectedIP, e.IP)
		}
	}
}

// Tests the Events creation from a Collection.
func TestNewEventCreation(t *testing.T) {
	t.Log("Creating a new Event object from Payload and Collection")
	env := payload.Environment_PRODUCTION
	parameter := payload.Parameter{
		Name:  proto.String("parameter"),
		Value: proto.String("test_parameter")}
	p := payload.Payload{
		At:         proto.Uint64(1454681104),
		Timezone:   proto.String("+02:00"),
		Event:      proto.String("Client.CreateUser"),
		Nr:         proto.Uint32(1),
		UserId:     proto.String("97421193"),
		TenantId:   proto.String("sdfghjkloiuytremiwoz"),
		Ip:         proto.String("214.160.227.22"),
		Country:    proto.String("UK"),
		Parameters: []*payload.Parameter{&parameter}}
	c := payload.Collection{
		DeviceId:        proto.String("a73b1c37-2c24-4786-af7a-16de88fbe23a"),
		ClientId:        proto.String("bce44f67b2661fd445d469b525b04f68"),
		Session:         proto.String("0e350a2cc31648bb24ba61eb14be337c"),
		Env:             &env,
		SystemVersion:   proto.String("10.10"),
		ProductVersion:  proto.String("1.1.2"),
		DeviceMake:      proto.String("Iphone"),
		DeviceModel:     proto.String("Iphone 6"),
		System:          proto.String("OSX"),
		SystemLanguage:  proto.String("DE"),
		Browser:         proto.String("Mozilla"),
		BrowserVersion:  proto.String("10.01.11"),
		ProductGitHash:  proto.String("5416a5889392d509e3bafcf40f6388e83aab23e6"),
		ProductLanguage: proto.String("HU"),
		Payloads:        []*payload.Payload{&p}}
	e := NewEvent(&c, &p)

	if exp := ConvertIsoformat(1454681104); e.At != exp {
		t.Errorf("Expected At was %s but it was %s instead", exp, e.At)
	}
	if exp := "+02:00"; e.Timezone != exp {
		t.Errorf("Expected Timezone was %s but it was %s instead", exp, e.Timezone)
	}
	if exp := "Client.CreateUser"; e.Event != exp {
		t.Errorf("Expected Event was %s but it was %s instead", exp, e.Event)
	}
	if exp := uint32(1); e.Nr != exp {
		t.Errorf("Expected Nr was %s but it was %d instead", exp, e.Nr)
	}
	if exp := "97421193"; e.UserID != exp {
		t.Errorf("Expected UserID was %s but it was %d instead", exp, e.UserID)
	}
	if exp := "sdfghjkloiuytremiwoz"; e.TenantID != exp {
		t.Errorf("Expected TenantID was %s but it was %d instead", exp, e.TenantID)
	}
	if exp := "214.160.227.22"; e.IP != exp {
		t.Errorf("Expected IP was %s but it was %s instead", exp, e.IP)
	}
	if exp := "UK"; e.Country != exp {
		t.Errorf("Expected Country was %s but it was %s instead", exp, e.Country)
	}
	if exp := "{\"parameter\":\"test_parameter\"}"; e.Parameters != exp {
		t.Errorf("Expected Parameters was %s but it was %s instead", exp, e.Parameters)
	}
	if exp := "a73b1c37-2c24-4786-af7a-16de88fbe23a"; e.DeviceID != exp {
		t.Errorf("Expected DeviceID was %s but it was %s instead", exp, e.DeviceID)
	}
	if exp := "bce44f67b2661fd445d469b525b04f68"; e.ClientID != exp {
		t.Errorf("Expected ClientID was %s but it was %s instead", exp, e.ClientID)
	}
	if exp := "0e350a2cc31648bb24ba61eb14be337c"; e.Session != exp {
		t.Errorf("Expected Session was %s but it was %s instead", exp, e.Session)
	}
	if exp := "PRODUCTION"; e.Env != exp {
		t.Errorf("Expected Env was %s but it was %s instead", exp, e.Env)
	}
	if exp := "10.10"; e.SystemVersion != exp {
		t.Errorf("Expected SystemVersion was %s but it was %s instead", exp, e.SystemVersion)
	}
	if exp := "Iphone"; e.DeviceMake != exp {
		t.Errorf("Expected DeviceMake was %s but it was %s instead", exp, e.DeviceMake)
	}
	if exp := "Iphone 6"; e.DeviceModel != exp {
		t.Errorf("Expected DeviceModel was %s but it was %s instead", exp, e.DeviceModel)
	}
	if exp := "OSX"; e.System != exp {
		t.Errorf("Expected System was %s but it was %s instead", exp, e.System)
	}
	if exp := "DE"; e.SystemLanguage != exp {
		t.Errorf("Expected SystemLanguage was %s but it was %s instead", exp, e.SystemLanguage)
	}
	if exp := "Mozilla"; e.Browser != exp {
		t.Errorf("Expected Browser was %s but it was %s instead", exp, e.Browser)
	}
	if exp := "10.01.11"; e.BrowserVersion != exp {
		t.Errorf("Expected BrowserVersion was %s but it was %s instead", exp, e.BrowserVersion)
	}
	if exp := "5416a5889392d509e3bafcf40f6388e83aab23e6"; e.ProductGitHash != exp {
		t.Errorf("Expected ProductGitHash was %s but it was %s instead", exp, e.ProductGitHash)
	}
	if exp := "HU"; e.ProductLanguage != exp {
		t.Errorf("Expected ProductLanguage was %s but it was %s instead", exp, e.ProductLanguage)
	}
}
