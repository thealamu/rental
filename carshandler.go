package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//carsHandler serves path /cars
func carsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pubCars, err := db.listPublicCars()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(pubCars)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
