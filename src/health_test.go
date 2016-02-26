package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	t.Log("Testing health handler")
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	HealthHandler(resp, req)

	if code := resp.Code; code != http.StatusOK {
		t.Errorf("Expected call to be successul. Got %d instead", code)
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json Content-Type. Got %s", contentType)
	}
}
