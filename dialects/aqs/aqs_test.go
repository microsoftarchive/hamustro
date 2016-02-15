package aqs

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	t.Log("Filled out configuration")
	c := Config{Account: "a", AccessKey: "k", QueueName: "q"}
	if !c.IsValid() {
		t.Errorf("Filled out configuration should be okay.")
	}

	t.Log("Incomplete configuration tests")
	c = Config{AccessKey: "k", QueueName: "q"}
	if c.IsValid() {
		t.Errorf("Account is not defined, should be invalid.")
	}

	c = Config{Account: "a", QueueName: "q"}
	if c.IsValid() {
		t.Errorf("AccessKey is not defined, should be invalid")
	}

	c = Config{Account: "a", AccessKey: "k"}
	if c.IsValid() {
		t.Errorf("QueueName is not defined, should be invalid")
	}
}
