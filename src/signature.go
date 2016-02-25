package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/wunderlist/hamustro/src/payload"
	"io"
)

// Returns the request's signature
func GetSignature(body []byte, time string) string {
	bodyHash := md5.New()
	io.WriteString(bodyHash, string(body[:]))

	requestHash := sha256.New()
	io.WriteString(requestHash, time)
	io.WriteString(requestHash, "|")
	io.WriteString(requestHash, hex.EncodeToString(bodyHash.Sum(nil)))
	io.WriteString(requestHash, "|")
	io.WriteString(requestHash, config.SharedSecret)

	return base64.StdEncoding.EncodeToString(requestHash.Sum(nil))
}

// Returns the protobuf message's session
func GetSession(c *payload.Collection) string {
	session := md5.New()
	io.WriteString(session, c.GetDeviceId())
	io.WriteString(session, ":")
	io.WriteString(session, c.GetClientId())
	io.WriteString(session, ":")
	io.WriteString(session, c.GetSystemVersion())
	io.WriteString(session, ":")
	io.WriteString(session, c.GetProductVersion())
	return hex.EncodeToString(session.Sum(nil))
}
