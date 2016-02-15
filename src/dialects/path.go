package dialects

import (
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"
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

// Resolves custom names in the path
func ResolvePath(basePath string) string {
	newPath := basePath
	if strings.Contains(basePath, "{date}") {
		newPath = strings.Replace(basePath, "{date}", time.Now().UTC().Format("2006-01-02"), -1)
	}
	return newPath
}

// Get a random name for the blob
func GetRandomPath(basePath string, extension string) string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fileName := fmt.Sprintf("%s-%s.%s.gz", timestamp, RandStringBytes(20), extension)
	return path.Join(ResolvePath(basePath), fileName)
}
