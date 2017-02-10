package dialects

import (
	"testing"
	"time"
)

// Generates a random string to usa as a part of a filename
func TestFunctionRandStringBytes(t *testing.T) {
	t.Log("Generates 10 length random string.")
	if s := RandStringBytes(10); len(s) != 10 {
		t.Errorf("Expected length was 10 but it was %d instead.", len(s))
	}
}

// Replacing the predefined placeholders in the filepath
func TestFunctionResolvePath(t *testing.T) {
	t.Log("Detecting placeholders in the file path")

	cases := []struct {
		Path     string
		Expected string
	}{
		{"directory/subdirectory", "directory/subdirectory"},
		{"directory/{date}", "directory/" + time.Now().UTC().Format("2006-01-02")},
		{"directory/{year}", "directory/" + time.Now().UTC().Format("2006")},
		{"directory/{month}", "directory/" + time.Now().UTC().Format("01")},
		{"directory/{day}", "directory/" + time.Now().UTC().Format("02")},
		{"directory/{hour}", "directory/" + time.Now().UTC().Format("15")},
		{"directory/{minute}", "directory/" + time.Now().UTC().Format("04")},
		{"directory/{second}", "directory/" + time.Now().UTC().Format("05")},
		{"directory/{year}/{month}/{day}", "directory/" + time.Now().UTC().Format("2006/01/02")}}

	for _, c := range cases {
		if p := ResolvePath(c.Path); p != c.Expected {
			t.Errorf("Expected path was %s but it was %s instead", c.Expected, p)
		}
	}
}

// Generates a random filepath for a single file
func TestFunctionGetRandomPath(t *testing.T) {
	t.Log("Generates random file name, for `csv` extension")

	p := GetRandomPath("", "csv", true)
	if ext := p[len(p)-6:]; ext != "csv.gz" {
		t.Errorf("Expected extension was csv.gz but it was %s instead", ext)
	}

	p = GetRandomPath("", "csv", false)
	if ext := p[len(p)-3:]; ext != "csv" {
		t.Errorf("Expected extension was csv but it was %s instead", ext)
	}

	p = GetRandomPath("", "csv", true)
	if exp := 38; len(p) != exp {
		t.Errorf("Expected length was %d without path but it was %d instead", exp, len(p))
	}

	p = GetRandomPath("", "csv", true)
	if p[0] == '/' {
		t.Errorf("Expected first charater was not / without path but it was %s", p[0])
	}

	p = GetRandomPath("dir", "csv", true)
	if exp := 42; len(p) != exp {
		t.Errorf("Expected length was %d with path but it was %d instead", exp, len(p))
	}

	p = GetRandomPath("dir/to/path/", "csv", true)
	if dir := p[0:12]; dir != "dir/to/path/" {
		t.Errorf("Expected directory path was dir/to/path/ but it was %s instead", dir)
	}

	p = GetRandomPath("dir/to/path", "csv", true)
	if dir := p[0:12]; dir != "dir/to/path/" {
		t.Errorf("Expected directory path was dir/to/path/ but it was %s instead", dir)
	}
}
