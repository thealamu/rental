package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseError(t *testing.T) {
	testErr := fmt.Errorf("Some Error")
	rRecorder := httptest.NewRecorder()
	respondError("Some Tag", rRecorder, 34, testErr.Error(), http.StatusTeapot)

	if rRecorder.Code != http.StatusTeapot {
		t.Errorf("respondError does not return the right status code, want %v, got %v", http.StatusTeapot, rRecorder.Code)
	}

	var errResp errorResponse
	err := json.NewDecoder(rRecorder.Body).Decode(&errResp)
	if err != nil {
		t.Error(err)
	}
	if errResp.Msg != "Some Error" {
		t.Errorf("respondError does not return the right message, want %v, got %v", "Some Error", errResp.Msg)
	}
	if errResp.Code != 34 {
		t.Errorf("respondError does not return the right error code, want %v, got %v", 34, errResp.Code)
	}
}

func TestRespondJSON(t *testing.T) {
	someData := struct {
		Data string
	}{
		Data: "Some String",
	}

	someDataBytes, err := json.Marshal(someData)
	if err != nil {
		t.Error(err)
	}

	rRecorder := httptest.NewRecorder()

	respondJSON(rRecorder, someData)

	respBodyStr := strings.TrimSpace(rRecorder.Body.String())
	someDataStr := strings.TrimSpace(string(someDataBytes))

	if respBodyStr != someDataStr {
		t.Errorf("respondJSON does not write passed in data, expected %v, got %v", someDataStr, respBodyStr)
	}

	//Test 500
	errRecorder := httptest.NewRecorder()
	respondJSON(errRecorder, func() {})

	if errRecorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusInternalServerError, errRecorder.Code)
	}
}
