package main

import (
	"context"
	"fmt"
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

func newAuthenticator() (*authenticator, error) {
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
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  authRedirectURL,
		Endpoint:     prov.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
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
