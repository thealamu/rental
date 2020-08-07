package main

import (
	"log"
	"net/http"
)

func getMerchantMeCars(w http.ResponseWriter, r *http.Request) {

}

func getMerchantMe(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchantme"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeSessionDB, "", http.StatusInternalServerError)
		return
	}

	//Get merchant details from db
	mcht, err := db.getMerchantForEmail(email)
	if err != nil {
		if err == errNotFound {
			//user is not a merchant
			log.Printf("%s: %v is not a merchant", tag, email)
			respondError(tag, w, failCodeAuth, "Not a merchant", http.StatusForbidden)
			return
		}
		log.Printf("%s: %v", tag, err)
		return
	}

	respondJSON(w, mcht)
}
