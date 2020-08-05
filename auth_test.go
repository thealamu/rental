package main

import (
	"os"
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	//test no domain
	_, err := newAuthenticator()
	if err != errNoDomain {
		t.Errorf("newAuthenticator: expected no domain error, got nil")
	}

	//test failing provider
	os.Setenv(domainKey, "https://someDomain")
	_, err = newAuthenticator()
	if err == nil {
		t.Errorf("newAuthenticator: expected invalid domain error, got nil")
	}

	os.Setenv(domainKey, "https://test-rental.us.auth0.com/")
	authen, err := newAuthenticator()
	if err != nil {
		t.Error(err)
	}
	if authen == nil {
		t.Errorf("newAuthenticator: expected an authenticator, got %v", authen)
	}
}
