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

	profInterface, ok := session.Values["profile"]
	if !ok {
		respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userProfile := profInterface.(map[string]interface{})
	username, ok := userProfile["name"]
	if !ok {
		respondError(tag, w, failCodeAuth, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Get merchant details from db
	mcht, err := db.getMerchantForName(username.(string))

	respondJSON(w, mcht)
}
