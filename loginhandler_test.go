package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

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
