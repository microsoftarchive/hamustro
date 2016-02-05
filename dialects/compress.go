package dialects

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generates an `n` length random string.
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Get a random name for the blob
func GetRandomPath(basePath string, extension string) string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := fmt.Sprintf("%s-%s.%s.gz", timestamp, RandStringBytes(20), extension)
	return path.Join(basePath, fileName)
}

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
