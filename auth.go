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

func getHost(r *http.Request) (string, error) {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	host, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		return "", err
	}

	return host.String(), nil
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

	clientID := os.Getenv(clientIDKey)

	host, err := getHost(r)
	if err != nil {
		return nil, err
	}

	parameters.Add("returnTo", host)
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

	host, err := getHost(r)
	if err != nil {
		return nil, err
	}

	authRedirectURL := host + authRedirectPath

	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  authRedirectURL,
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
