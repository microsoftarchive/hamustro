package dialects

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	t.Log("Generates 10 length random string.")
	if s := RandStringBytes(10); len(s) != 10 {
		t.Errorf("Expected length was 10 but it was %d instead.", len(s))
	}
}

func TestGetRandomPath(t *testing.T) {
	t.Log("Generates random file name, for `csv` extension")

	p := GetRandomPath("", "csv")
	if ext := p[len(p)-6:]; ext != "csv.gz" {
		t.Errorf("Expected extension was csv.gz but it was %s instead", ext)
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
