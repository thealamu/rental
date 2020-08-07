package main

import (
	"net/http"
	"os"
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	testRequest, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}

	//test no domain
	_, err = newAuthenticator(testRequest)
	if err != errNoDomain {
		t.Errorf("newAuthenticator: expected no domain error, got nil")
	}

	//test failing provider
	os.Setenv(domainKey, "https://someDomain")
	_, err = newAuthenticator(testRequest)
	if err == nil {
		t.Errorf("newAuthenticator: expected invalid domain error, got nil")
	}

	os.Setenv(domainKey, "https://test-rental.us.auth0.com/")
	authen, err := newAuthenticator(testRequest)
	if err != nil {
		t.Error(err)
	}
	if authen == nil {
		t.Errorf("newAuthenticator: expected an authenticator, got %v", authen)
	}
}
