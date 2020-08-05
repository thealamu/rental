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
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
		return
	}

	paramMerchant := mux.Vars(r)["merchant"]

	mcht, err := db.getMiniMerchantForName(paramMerchant)
	if err != nil {
		rspErrCode := http.StatusInternalServerError
		if err == errNotFound {
			rspErrCode = http.StatusNotFound
		}
		respondError(tag, w, failCodeBadParameter, err.Error(), rspErrCode)
		return
	}

	respondJSON(w, mcht)
}
