package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getSingleMiniMerchant(w http.ResponseWriter, r *http.Request) {
	tag := "handler.merchant"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	paramMerchant := mux.Vars(r)["merchant"]

	mcht, err := db.getMiniMerchantForName(paramMerchant)
	if err != nil {
		log.Printf("%s: %v for merchant param '%s'", tag, err, paramMerchant)

		rspErr := http.StatusInternalServerError
		if err == errNotFound {
			rspErr = http.StatusNotFound
		}

		w.WriteHeader(rspErr)
		return
	}

	respondJSON(tag, w, mcht)
}
