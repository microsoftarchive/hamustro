package dialects

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
	"time"
)

func TestRandStringBytes(t *testing.T) {
	t.Log("Generates 10 length random string.")
	if s := RandStringBytes(10); len(s) != 10 {
		t.Errorf("Expected length was 10 but it was %d instead.", len(s))
	}
}

func TestResolvePath(t *testing.T) {
	t.Log("Detecting placeholders in the file path")

	p := ResolvePath("directory/subdirectory")
	if exp := "directory/subdirectory"; exp != p {
		t.Errorf("Expected path was %s but it was %s instead", exp, p)
	}

	p = ResolvePath("directory/{date}")
	if exp := "directory/" + time.Now().UTC().Format("2006-01-02"); exp != p {
		t.Errorf("Expected path was %s but it was %s instead", exp, p)
	}
}

func TestGetRandomPath(t *testing.T) {
	t.Log("Generates random file name, for `csv` extension")

	p := GetRandomPath("", "csv")
	if ext := p[len(p)-6:]; ext != "csv.gz" {
		t.Errorf("Expected extension was csv.gz but it was %s instead", ext)
	}

	p = GetRandomPath("", "csv")
	if exp := 38; len(p) != exp {
		t.Errorf("Expected length was %d without path but it was %d instead", exp, len(p))
	}

	p = GetRandomPath("dir", "csv")
	if exp := 42; len(p) != exp {
		t.Errorf("Expected length was %d with path but it was %d instead", exp, len(p))
	}

	p = GetRandomPath("dir/to/path/", "csv")
	if dir := p[0:12]; dir != "dir/to/path/" {
		t.Errorf("Expected directory path was dir/to/path/ but it was %s instead", dir)
	}

	p = GetRandomPath("dir/to/path", "csv")
	if dir := p[0:12]; dir != "dir/to/path/" {
		t.Errorf("Expected directory path was dir/to/path/ but it was %s instead", dir)
	}
}

func TestCompess(t *testing.T) {
	t.Log("Compressing the string `hamustro`")
	b := new(bytes.Buffer)
	b.Write([]byte("hamustro"))

	cb, err := Compress(b)
	if err != nil {
		t.Errorf("Compress is failed: %s", err.Error())
	}

	gz, err := gzip.NewReader(cb)
	if err != nil {
		t.Errorf("Decompressing is failed: %s", err.Error())
	}
	defer gz.Close()

	fb := new(bytes.Buffer)
	io.Copy(fb, gz)
	if fb.String() != "hamustro" {
		t.Error("Compress function is not generating valid gzip archive")
	}
}
