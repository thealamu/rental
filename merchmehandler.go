package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func createMerchantMeCar(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchantmecreatecar"

	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	var carItem car
	err = json.NewDecoder(r.Body).Decode(&carItem)
	if err != nil {
		respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
		return
	}

	err = verifyCarItem(&carItem)
	if err != nil {
		respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.createMerchantCar(&carItem)
	if err != nil {
		respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, carItem)
}

func getMerchantMeCars(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchantmecars"

	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
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

	respondJSON(w, db.getAuthedMerchantCars(mcht.Name))
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
