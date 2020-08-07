package main

import (
	"log"
	"net/http"
)

//Handles path /cars/search
func handleCarsSearch(w http.ResponseWriter, r *http.Request) {
	tag := "handler.searchcars"

	q := r.URL.Query().Get("q")
	if q == "" {
		respondError(tag, w, failCodeBadParameter, "No search query", http.StatusBadRequest)
		return
	}

	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
		return
	}

	respondJSON(w, db.searchPublicCars(q))
}
