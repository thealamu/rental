package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//carHandler serves path /cars/{car_id}
func getSinglePublicCar(w http.ResponseWriter, r *http.Request) {
	tag := "handler.car"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		respondError(tag, w, failCodeDB, "", http.StatusInternalServerError)
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

	pubCar, err := db.getPublicCarForID(uint(paramCarID))
	if err != nil {
		rspErr := http.StatusInternalServerError
		if err == errNotFound {
			rspErr = http.StatusNotFound
		}
		respondError(tag, w, failCodeBadParameter, err.Error(), rspErr)
		return
	}

	respondJSON(w, pubCar)
}

//carsHandler serves path /cars
func getPublicCars(w http.ResponseWriter, r *http.Request) {
	tag := "handler.cars"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	pubCars, err := db.listPublicCars()
	if err != nil {
		respondError(tag, w, failCodeDB, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, pubCars)
}
