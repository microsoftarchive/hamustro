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
	rules := []struct {
		Keyword string
		Format  string
	}{
		{"{date}", "2006-01-02"},
		{"{year}", "2006"},
		{"{month}", "01"},
		{"{day}", "02"},
		{"{hour}", "15"},
		{"{minute}", "04"},
		{"{second}", "05"}}

	newPath := basePath
	for _, rule := range rules {
		if strings.Contains(newPath, rule.Keyword) {
			newPath = strings.Replace(newPath, rule.Keyword, time.Now().UTC().Format(rule.Format), -1)
		}
	}
	return newPath
}

// Get a random name for the blob
func GetRandomPath(basePath string, extension string, compress bool) string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	compressedExtension := ""
	if compress {
		compressedExtension = ".gz"
	}
	fileName := fmt.Sprintf("%s-%s.%s%s", timestamp, RandStringBytes(20),
		extension, compressedExtension)
	return path.Join(ResolvePath(basePath), fileName)
}
