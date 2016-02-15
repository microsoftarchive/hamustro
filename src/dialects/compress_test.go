package dialects

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestCompress(t *testing.T) {
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
