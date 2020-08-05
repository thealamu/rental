package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

func TestAppendRedirURL(t *testing.T) {
	testState := "somestatedata"
	testRequest, _ := http.NewRequest(http.MethodGet, "/auth/login/callback?state_url=somestateurl", nil)

	appended, err := appendRedirURL(testState, testRequest)
	if err != nil {
		t.Errorf("appendRedirURL expected no error, got %v", err)
	}

	if appended != "somestatedata?state_url=somestateurl" {
		t.Errorf("appendRedirURL returns wrong state string")
	}

	badRequest, _ := http.NewRequest(http.MethodGet, "/auth/login/callback", nil)
	_, err = appendRedirURL(testState, badRequest)
	if err == nil {
		t.Errorf("appendRedirURL: expected an error, got nil")
	}
}

func TestGetRedirURL(t *testing.T) {
	if getRedirURL("some?state_url=data") != "data" {
		t.Errorf("getRedirURL does not return correct state_url")
	}
}

func TestHandleLogin(t *testing.T) {
	os.Setenv(domainKey, "https://test-rental.us.auth0.com/")
	testRequest, err := http.NewRequest(http.MethodGet, "/auth/login?state_url=someurl", nil)
	if err != nil {
		t.Error(err)
	}

	os.Setenv("RTL_STOREKEY", "somekey")
	initSessionStore()

	rRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/auth/login", handleLogin)

	router.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusTemporaryRedirect, rRecorder.Code)
	}

	//test no state_url
	noStateRequest, err := http.NewRequest(http.MethodGet, "/auth/login", nil)
	if err != nil {
		t.Error(err)
	}

	errRecorder := httptest.NewRecorder()
	router.ServeHTTP(errRecorder, noStateRequest)

	if errRecorder.Code != http.StatusBadRequest {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusBadRequest, errRecorder.Code)
	}

	//test bad authenticator
	os.Setenv(domainKey, "baddomain.bad")
	badDomainRequest, err := http.NewRequest(http.MethodGet, "/auth/login?state_url=someurl", nil)
	if err != nil {
		t.Error(err)
	}
	errRecorder = httptest.NewRecorder()
	router.ServeHTTP(errRecorder, badDomainRequest)

	if errRecorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusInternalServerError, errRecorder.Code)
	}
}
