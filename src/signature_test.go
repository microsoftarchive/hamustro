package main

import (
	"strconv"
	"testing"
)

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
	collection := GetTestPayloadCollection(97421193, 1)
	if exp := "244f056dee6d475ec673ea0d20b69bab"; GetSession(collection) != exp {
		t.Errorf("Expected session was %s and it was %s instead", exp, GetSession(collection))
	}
}
