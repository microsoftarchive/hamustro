package s3

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	t.Log("Filled out configuration")
	c := Config{AccessKeyID: "a", SecretAccessKey: "k", Bucket: "b", BlobPath: "bp", Endpoint: "ep", Region: "r", FileFormat: "ff"}
	if !c.IsValid() {
		t.Errorf("Filled out configuration should be okay.")
	}

	t.Log("Incomplete configuration tests")
	c = Config{SecretAccessKey: "k", Bucket: "b", BlobPath: "bp", Endpoint: "ep", Region: "r", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("AccessKeyID is not defined, should be invalid.")
	}

	c = Config{AccessKeyID: "a", Bucket: "b", BlobPath: "bp", Endpoint: "ep", Region: "r", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("SecretAccessKey is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", BlobPath: "bp", Endpoint: "ep", Region: "r", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("Bucket is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", Bucket: "b", Endpoint: "ep", Region: "r", FileFormat: "ff"}
	if !c.IsValid() {
		t.Errorf("BlobPath is not defined, should be empty and still valid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", Bucket: "b", BlobPath: "bp", Region: "r", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("Endpoint is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", Bucket: "b", BlobPath: "bp", Endpoint: "ep", FileFormat: "ff"}
	if c.IsValid() {
		t.Errorf("Region is not defined, should be invalid")
	}

	c = Config{AccessKeyID: "a", SecretAccessKey: "k", Bucket: "b", BlobPath: "bp", Endpoint: "ep", Region: "r"}
	if c.IsValid() {
		t.Errorf("FileFormat is not defined, should be invalid")
	}
}
