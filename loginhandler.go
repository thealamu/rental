package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

var authRedirectURL = "http://localhost:8080/auth/login/callback"

var (
	errNoAcctName  = fmt.Errorf("Account name not set for merchant")
	errNoStateURL  = fmt.Errorf("No state_url set")
	errNoLoginType = fmt.Errorf("No login_type set")
)

const (
	acctTypeCustomer = "customer"
	acctTypeMerchant = "merchant"
)

const (
	loginTypeLogin  = "login"
	loginTypeSignup = "signup"
)

func getRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(b)
	return state, nil
}

func saveLoginParams(state string, w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		return err
	}

	session.Values["state"] = state

	loginType := r.URL.Query().Get("login_type")
	if loginType == "" {
		return errNoLoginType
	}

	stateURL := r.URL.Query().Get("state_url")
	if stateURL == "" {
		return errNoStateURL
	}
	session.Values["state_url"] = stateURL

	if loginType == loginTypeSignup {
		acctType := r.URL.Query().Get("account_type")
		if acctType == "" {
			//default account type is customer
			acctType = acctTypeCustomer
		}

		acctName := r.URL.Query().Get("account_name")
		if acctType == acctTypeMerchant && acctName == "" {
			//account name is required for a merchant
			return errNoAcctName
		}

		session.Values["account_type"] = acctType
		session.Values["account_name"] = acctName
	}

	session.Values["login_type"] = loginType

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	tag := "handler.login"

	state, err := getRandomState()
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeUnknown, "", http.StatusInternalServerError)
		return
	}

	err = saveLoginParams(state, w, r)
	if err != nil {
		if err == errNoAcctName || err == errNoStateURL || err == errNoLoginType {
			respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "Session failure", http.StatusInternalServerError)
		return
	}

	authen, err := newAuthenticator()
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeAuth, "Auth configuration failure", http.StatusInternalServerError)
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
		respondError(tag, w, failCodeAuth, "Auth configuration failure", http.StatusInternalServerError)
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
	fmt.Println(profile)

	stateURL, ok := session.Values["state_url"].(string)
	if !ok {
		log.Printf("%s: state_url not set in session", tag)
		return
	}

	loginType, ok := session.Values["login_type"].(string)
	if !ok {
		log.Printf("%s: login_type not set in sesion", tag)
		return
	}

	if loginType == loginTypeLogin {
		http.Redirect(w, r, stateURL, http.StatusTemporaryRedirect)
		return
	}

	//trying to signup
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
		return
	}

	acctType, ok := session.Values["account_type"].(string)
	if !ok {
		log.Printf("%s: account_type not set in session", tag)
		return
	}

	if acctType == acctTypeCustomer {
		//TODO: Do stuff
	} else if acctType == acctTypeMerchant {
		acctName, ok := session.Values["account_name"].(string)
		if !ok {
			log.Printf("%s: account_name not set in session", tag)
			return
		}
		email, err := getProfileValue(r, "name")
		if err != nil {
			log.Printf("%s: %v", tag, err)
			return
		}
		db.createMerchant(acctName, email)
	} else {
		log.Printf("%s: account_type %s not known", tag, acctType)
	}

	http.Redirect(w, r, stateURL, http.StatusTemporaryRedirect)
}
