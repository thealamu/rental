package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRootHandler(t *testing.T) {
	testRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", getCommonEndpoints)
	rRecorder := httptest.NewRecorder()

	r.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusOK {
		t.Errorf("Bad status code, got %v, expected %v", rRecorder.Code, http.StatusOK)
	}

	data, err := json.Marshal(commonEndpoints)
	if err != nil {
		t.Error(err)
	}
	expected := string(data)
	respBody := rRecorder.Body.String()

	if respBody != expected {
		t.Errorf("Wrong response body")
	}
}
