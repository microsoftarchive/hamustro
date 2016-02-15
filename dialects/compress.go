package dialects

import (
	"bytes"
	"compress/gzip"
)

// Compress the given string
func Compress(msg *bytes.Buffer) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	gz := gzip.NewWriter(b)
	if _, err := gz.Write(msg.Bytes()); err != nil {
		return b, err
	}
	if err := gz.Flush(); err != nil {
		return b, err
	}
	if err := gz.Close(); err != nil {
		return b, err
	}
	return b, nil
}
