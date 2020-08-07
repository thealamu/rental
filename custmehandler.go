package main

import (
	"log"
	"net/http"
)

func getCustomerMe(w http.ResponseWriter, r *http.Request) {
	tag := "handler.custme"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get customer details from db
	cust, err := db.getCustomerForEmail(email)
	if err != nil {
		if err == errNotFound {
			//user is not a customer
			respondError(tag, w, failCodeAuth, "Not a customer", http.StatusForbidden)
			return
		}
		log.Printf("%s: %v", tag, err)
		return
	}

	respondJSON(w, cust)
}
