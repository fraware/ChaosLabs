package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartExperimentHandler(t *testing.T) {
	// Create a sample experiment request with scheduling in the near future.
	expReq := ExperimentRequest{
		Name:           "Test Experiment",
		Description:    "Testing scheduling logic",
		ExperimentType: "cpu-stress",
		Duration:       10,
		StartTime:      time.Now().Add(2 * time.Second),
	}
	body, err := json.Marshal(expReq)
	if err != nil {
		t.Fatalf("Failed to marshal experiment request: %v", err)
	}

	req, err := http.NewRequest("POST", "/start", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Call the handler.
	handler := http.HandlerFunc(startExperimentHandler)
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
	if resp["status"] != "scheduled" {
		t.Errorf("Expected response status 'scheduled', got '%s'", resp["status"])
	}
}
