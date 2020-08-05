package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var authRedirectURL = "http://localhost:8080/auth/login/callback"

func handleLogin(w http.ResponseWriter, r *http.Request) {
	tag := "handler.login"
	//set the state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeUnknown, "", http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)
	state, err = appendRedirURL(state, r)
	if err != nil {
		respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "auth-session")
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
		return
	}

	session.Values["state"] = state
	if err := session.Save(r, w); err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
		return
	}

	authen, err := newAuthenticator()
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Auth config failure", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authen.oauthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	tag := "handler.authCallback"
	session, err := store.Get(r, "auth-session")
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
		return
	}

	stateQuery := r.URL.Query().Get("state")
	if stateQuery != session.Values["state"] {
		respondError(tag, w, failCodeBadParameter, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	authen, err := newAuthenticator()
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Auth config failure", http.StatusInternalServerError)
		return
	}

	codeQuery := r.URL.Query().Get("code")
	token, err := authen.oauthConfig.Exchange(r.Context(), codeQuery)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("%s: No id_token field in oauth2 token", tag)
		respondError(tag, w, failCodeAuth, "Auth Failure", http.StatusInternalServerError)
		return
	}

	idToken, err := authen.provider.Verifier(authen.oidcConfig).Verify(r.Context(), rawIDToken)
	if err != nil {
		log.Printf("%s: Failed to verify ID Token", tag)
		respondError(tag, w, failCodeAuth, "Auth Failure", http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Auth Failure", http.StatusInternalServerError)
		return
	}

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	if err := session.Save(r, w); err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Auth Failure", http.StatusInternalServerError)
		return
	}

	//go to the redirect url inside state
	http.Redirect(w, r, getRedirURL(stateQuery), http.StatusTemporaryRedirect)
}

func getRedirURL(state string) string {
	return strings.Split(state, "?state_url=")[1]
}

func appendRedirURL(state string, r *http.Request) (string, error) {
	redirURL := r.URL.Query().Get("state_url")
	if redirURL == "" {
		return "", fmt.Errorf("state_url not set")
	}
	return state + "?state_url=" + redirURL, nil
}
