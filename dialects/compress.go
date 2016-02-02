package dialects

import (
	"bytes"
	"compress/gzip"
	"math/rand"
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
func GetRandomPath(path string) string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := timestamp + "-" + RandStringBytes(20) + ".json.gz"
	filePath := path + fileName
	return filePath
}

// Compress the given string
func Compress(msg *string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(*msg)); err != nil {
		return &b, err
	}
	if err := gz.Flush(); err != nil {
		return &b, err
	}
	if err := gz.Close(); err != nil {
		return &b, err
	}
	return &b, nil
}
