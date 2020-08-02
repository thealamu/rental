package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//carsHandler serves path /cars
func carsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := newDatabase(nil)
	if err != nil {
		log.Println(err)
		return
	}

	pubCars, err := db.listPublicCars()
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(pubCars)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(data)
}
