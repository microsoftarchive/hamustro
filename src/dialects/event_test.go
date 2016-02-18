package dialects

import (
	"../payload"
	"github.com/golang/protobuf/proto"
	"reflect"
	"testing"
)

// Returns an Event for testing purposes
func GetTestEvent(userId uint32) *Event {
	return &Event{
		DeviceID:       "a73b1c37-2c24-4786-af7a-16de88fbe23a",
		ClientID:       "bce44f67b2661fd445d469b525b04f68",
		Session:        "244f056dee6d475ec673ea0d20b69bab",
		Nr:             1,
		SystemVersion:  "10.10",
		ProductVersion: "1.1.2",
		At:             "2016-02-05T15:05:04",
		Event:          "Client.CreateUser",
		System:         "OSX",
		ProductGitHash: "5416a5889392d509e3bafcf40f6388e83aab23e6",
		UserID:         userId,
		IP:             "214.160.227.22",
		Parameters:     "",
		IsTesting:      false}
}

// Converts and Event into a list of string
func TestEventStringConversion(t *testing.T) {
	t.Log("Converting Event into a list of strings")
	e := GetTestEvent(97421193)
	exp := []string{
		"a73b1c37-2c24-4786-af7a-16de88fbe23a",
		"bce44f67b2661fd445d469b525b04f68",
		"244f056dee6d475ec673ea0d20b69bab",
		"1",
		"10.10",
		"1.1.2",
		"2016-02-05T15:05:04",
		"Client.CreateUser",
		"OSX",
		"5416a5889392d509e3bafcf40f6388e83aab23e6",
		"97421193",
		"214.160.227.22",
		"",
		"false"}
	if !reflect.DeepEqual(e.String(), exp) {
		t.Error("Expected event's string is not matched")
	}
}

// Tests the Events creation from a Collection.
func TestNewEventCreation(t *testing.T) {
	t.Log("Creating a new Event object from Payload and Collection")
	p := payload.Payload{
		At:         proto.Uint64(1454681104),
		Event:      proto.String("Client.CreateUser"),
		Nr:         proto.Uint32(1),
		UserId:     proto.Uint32(97421193),
		Ip:         proto.String("214.160.227.22"),
		Parameters: proto.String(""),
		IsTesting:  proto.Bool(false)}
	c := payload.Collection{
		DeviceId:       proto.String("a73b1c37-2c24-4786-af7a-16de88fbe23a"),
		ClientId:       proto.String("bce44f67b2661fd445d469b525b04f68"),
		Session:        proto.String("244f056dee6d475ec673ea0d20b69bab"),
		SystemVersion:  proto.String("10.10"),
		ProductVersion: proto.String("1.1.2"),
		System:         proto.String("OSX"),
		ProductGitHash: proto.String("5416a5889392d509e3bafcf40f6388e83aab23e6"),
		Payloads:       []*payload.Payload{&p}}
	e := NewEvent(&c, &p)

	if exp := "a73b1c37-2c24-4786-af7a-16de88fbe23a"; e.DeviceID != exp {
		t.Errorf("Expected DeviceID was %s but it was %s instead", exp, e.DeviceID)
	}
	if exp := "bce44f67b2661fd445d469b525b04f68"; e.ClientID != exp {
		t.Errorf("Expected ClientID was %s but it was %s instead", exp, e.ClientID)
	}
	if exp := "244f056dee6d475ec673ea0d20b69bab"; e.Session != exp {
		t.Errorf("Expected Session was %s but it was %s instead", exp, e.Session)
	}
	if exp := uint32(1); e.Nr != exp {
		t.Errorf("Expected Nr was %s but it was %d instead", exp, e.Nr)
	}
	if exp := "10.10"; e.SystemVersion != exp {
		t.Errorf("Expected SystemVersion was %s but it was %s instead", exp, e.SystemVersion)
	}
	if exp := "1.1.2"; e.ProductVersion != exp {
		t.Errorf("Expected ProductVersion was %s but it was %s instead", exp, e.ProductVersion)
	}
	if exp := ConvertIsoformat(1454681104); e.At != exp {
		t.Errorf("Expected At was %s but it was %s instead", exp, e.At)
	}
	if exp := "Client.CreateUser"; e.Event != exp {
		t.Errorf("Expected Event was %s but it was %s instead", exp, e.Event)
	}
	if exp := "OSX"; e.System != exp {
		t.Errorf("Expected System was %s but it was %s instead", exp, e.System)
	}
	if exp := "5416a5889392d509e3bafcf40f6388e83aab23e6"; e.ProductGitHash != exp {
		t.Errorf("Expected ProductGitHash was %s but it was %s instead", exp, e.ProductGitHash)
	}
	if exp := uint32(97421193); e.UserID != exp {
		t.Errorf("Expected UserID was %s but it was %d instead", exp, e.UserID)
	}
	if exp := "214.160.227.22"; e.IP != exp {
		t.Errorf("Expected IP was %s but it was %s instead", exp, e.IP)
	}
	if exp := ""; e.Parameters != exp {
		t.Errorf("Expected Parameters was %s but it was %s instead", exp, e.Parameters)
	}
	if exp := false; e.IsTesting != exp {
		t.Errorf("Expected IsTesting was %s but it was %s instead", exp, e.IsTesting)
	}
}
