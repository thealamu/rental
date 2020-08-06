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

	username, err := getProfileValue(r, "name")
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "", http.StatusInternalServerError)
		return
	}

	//Get merchant details from db
	mcht, err := db.getMerchantForName(username)
	if err != nil {
		if err == errNotFound {
			//user is not a merchant
			log.Printf("%s: %v is not a merchant", tag, username)
			respondError(tag, w, failCodeAuth, "Not a merchant", http.StatusForbidden)
			return
		}
		log.Printf("%s: %v", tag, err)
		return
	}

	respondJSON(w, mcht)
}
