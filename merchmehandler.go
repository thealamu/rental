package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getMerchantMeSingleCar(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchantcar"

	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	param := mux.Vars(r)["car_id"]
	paramCarID, err := strconv.Atoi(param)
	if err != nil {
		//err should usually be nil because the router enforces the constraints
		//for a car id.
		//bad car id
		respondError(tag, w, failCodeBadParameter, err.Error(), http.StatusBadRequest)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
	}

	mcht, err := db.getMerchantForEmail(email)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
	}

	mchtCar, err := db.getMerchantCarForID(mcht.Name, uint(paramCarID))
	if err != nil {
		rspErr := http.StatusInternalServerError
		if err == errNotFound {
			rspErr = http.StatusNotFound
		}
		respondError(tag, w, failCodeBadParameter, err.Error(), rspErr)
		return
	}

	respondJSON(w, mchtCar)
}

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

	carItem.ID = db.getNextCarID()

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

	host, err := getHost(r)
	if err != nil {
		respondError(tag, w, failCodeUnknown, err.Error(), http.StatusInternalServerError)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.updateMerchantCarCount(email, carItem.IsPublic)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s/merchants/me/cars/%d", host, carItem.ID))
	w.WriteHeader(http.StatusCreated)
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
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get merchant details from db
	mcht, err := db.getMerchantForEmail(email)
	if err != nil {
		if err == errNotFound {
			//user is not a merchant
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
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	email, err := getProfileValue(r, "email")
	if err != nil {
		respondError(tag, w, failCodeSessionDB, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get merchant details from db
	mcht, err := db.getMerchantForEmail(email)
	if err != nil {
		if err == errNotFound {
			//user is not a merchant
			respondError(tag, w, failCodeAuth, "Not a merchant", http.StatusForbidden)
			return
		}
		log.Printf("%s: %v", tag, err)
		return
	}

	respondJSON(w, mcht)
}
