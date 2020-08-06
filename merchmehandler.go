package main

import (
	"log"
	"net/http"
)

func getMerchantMe(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchantme"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
		return
	}

	//get authed user
	session, err := store.Get(r, "auth-session")
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "Session store failure", http.StatusInternalServerError)
		return
	}

	profInterface, _ := session.Values["profile"]
	usernameInterface, ok := profInterface.(map[string]interface{})["name"]
	username, ok := usernameInterface.(string)
	if !ok {
		log.Printf("%s: %s for %v", tag, "name field not defined in profile", usernameInterface)
		respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("%s: %s is logged in", tag, username)

	//Get merchant details from db
	mcht, err := db.getMerchantForName(username)
	if mcht.Name == "" {
		//user is not a merchant
		log.Printf("%s: %v is not a merchant", tag, username)
		respondError(tag, w, failCodeAuth, "Not a merchant", http.StatusForbidden)
		return
	}

	respondJSON(w, mcht)
}
