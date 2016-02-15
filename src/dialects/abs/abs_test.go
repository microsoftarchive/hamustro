package abs

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	t.Log("Filled out configuration")
	c := Config{Account: "a", AccessKey: "k", Container: "c", BlobPath: "bp", FileFormat: "ff"}
	if !c.IsValid() {
		t.Errorf("Filled out configuration should be okay.")
	}

	t.Log("Incomplete configuration tests")
	c = Config{AccessKey: "k", Container: "c", BlobPath: "bp", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("Account is not defined, should be invalid.")
	}

	c = Config{Account: "a", Container: "c", BlobPath: "bp", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("AccessKey is not defined, should be invalid")
	}

	c = Config{Account: "a", AccessKey: "k", BlobPath: "bp", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("Container is not defined, should be invalid")
	}

	c = Config{Account: "a", AccessKey: "k", Container: "c", FileFormat: "ff"}
	if !c.IsValid() {
		t.Errorf("BlobPath is not defined, should be empty and still valid")
	}

	c = Config{Account: "a", AccessKey: "k", Container: "c", BlobPath: "bp"}
	if c.IsValid() {
		t.Errorf("FileFormat is not defined, should be invalid")
	}
}
