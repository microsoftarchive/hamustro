package main

import (
	"./payload"
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
)

// Returns a Collection
func GetTestCollection(userId uint32) *payload.Collection {
	p := payload.Payload{
		At:         proto.Uint64(1454681104),
		Event:      proto.String("Client.CreateUser"),
		Nr:         proto.Uint32(1),
		UserId:     proto.Uint32(userId),
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
	return &c
}

// Generates a signature for a given string and an EPOCH timestamp
func TestFunctionGetSignature(t *testing.T) {
	t.Log("Generating signature for a string.")
	config = &Config{SharedSecret: "ultrasafesecret"}

	signature := GetSignature([]byte("something"), strconv.Itoa(1454514088))
	if exp := "DAfTAP+9T/K/N08k+nwRTWNpfacimS8DJcQG1I4+Moo="; exp != signature {
		t.Errorf("Expected signature was %s and it was %s instead.", exp, signature)
	}
}

// Generates a session identifier for a payload
func TestFunctionGetSession(t *testing.T) {
	t.Log("Generating a session for a payload's collection.")
	collection := GetTestCollection(97421193)
	if exp := "244f056dee6d475ec673ea0d20b69bab"; GetSession(collection) != exp {
		t.Errorf("Expected session was %s and it was %s instead", exp, GetSession(collection))
	}
}
