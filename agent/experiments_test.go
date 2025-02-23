package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInjectHandler(t *testing.T) {
	// Create a sample injection request.
	injReq := InjectionRequest{
		ExperimentType: "cpu-stress",
		Duration:       5,
		CPUWorkers:     2,
	}
	body, err := json.Marshal(injReq)
	if err != nil {
		t.Fatalf("Failed to marshal injection request: %v", err)
	}

	req, err := http.NewRequest("POST", "/inject", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Call the injection handler.
	handler := http.HandlerFunc(injectHandler)
	handler.ServeHTTP(rr, req)

	// Check status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	// Check response body.
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if resp["status"] != "injected" {
		t.Errorf("Expected response status 'injected', got '%s'", resp["status"])
	}
}
