package sns

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	t.Log("Filled out configuration")
	c := Config{AccessKeyID: "a", SecretAccessKey: "k", TopicArn: "t", Region: "r"}
	if !c.IsValid() {
		t.Errorf("Filled out configuration should be okay.")
	}

	t.Log("Incomplete configuration tests")
	c = Config{SecretAccessKey: "k", TopicArn: "t", Region: "r"}
	if c.IsValid() {
		t.Errorf("AccessKeyID is not defined, should be invalid.")
	}

	c = Config{AccessKeyID: "a", TopicArn: "t", Region: "r"}
	if c.IsValid() {
		t.Errorf("SecretAccessKey is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", Region: "r"}
	if c.IsValid() {
		t.Errorf("TopicArn is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", TopicArn: "t"}
	if c.IsValid() {
		t.Errorf("Region is not defined, should be invalid")
	}
}
