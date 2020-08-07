package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	domainKey    = "RTL_DOMAIN"
	clientIDKey  = "RTL_CLIENT_ID"
	clientSecKey = "RTL_CLIENT_SECRET"
)

var (
	errNoDomain       = fmt.Errorf("%s not set in environment", domainKey)
	errNoClientID     = fmt.Errorf("%s not set in environment", clientIDKey)
	errNoClientSecret = fmt.Errorf("%s not set in environment", clientSecKey)
)

type authenticator struct {
	provider    *oidc.Provider
	oidcConfig  *oidc.Config
	oauthConfig oauth2.Config
	ctx         context.Context
}

func newLogoutURL(r *http.Request) (*url.URL, error) {
	domain := os.Getenv(domainKey)
	if domain == "" {
		return nil, errNoDomain
	}

	logoutURL, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}

	logoutURL.Path += "v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		return nil, err
	}

	clientID := os.Getenv(clientIDKey)

	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", clientID)

	logoutURL.RawQuery = parameters.Encode()

	return logoutURL, nil
}

func newAuthenticator(r *http.Request) (*authenticator, error) {
	ctx := context.Background()

	domain := os.Getenv(domainKey)
	if domain == "" {
		return nil, errNoDomain
	}
	prov, err := oidc.NewProvider(ctx, domain)
	if err != nil {
		return nil, err
	}

	clientID := os.Getenv(clientIDKey)
	clientSecret := os.Getenv(clientSecKey)

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	authRedirectURL, err := url.Parse(scheme + "://" + r.Host + authRedirectPath)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  authRedirectURL.String(),
		Endpoint:     prov.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	return &authenticator{
		provider:    prov,
		oidcConfig:  oidcConfig,
		oauthConfig: conf,
		ctx:         ctx,
	}, nil
}
