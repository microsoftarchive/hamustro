package file

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	cases := []struct {
		Config  *Config
		IsValid bool
	}{
		{&Config{FilePath: "fp", FileFormat: "ff"}, true},
		{&Config{FilePath: "fp", FileFormat: "ff", Compress: true}, true},
		{&Config{FilePath: "fp"}, false},
		{&Config{FileFormat: "fp"}, false},
		{&Config{Compress: true}, false},
	}

	for _, c := range cases {
		if c.Config.IsValid() != c.IsValid {
			t.Errorf("Configuration is not valid")
		}
	}

	c := Config{FilePath: "fp", FileFormat: "ff"}
	if c.Compress != false {
		t.Errorf("Compression should be disabled as default")
	}
}
