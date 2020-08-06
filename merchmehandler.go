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

	username := getSessionUsername(r)

	//Get merchant details from db
	mcht, err := db.getMerchantForName(username)
	if err != nil {
		//user is not a merchant
		log.Printf("%s: %v is not a merchant", tag, username)
		respondError(tag, w, failCodeAuth, "Not a merchant", http.StatusForbidden)
		return
	}

	respondJSON(w, mcht)
}
