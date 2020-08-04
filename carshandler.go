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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	param := mux.Vars(r)["car_id"]
	paramCarID, err := strconv.Atoi(param)
	if err != nil {
		//err should usually be nil because the router enforces the constraints
		//for a car id.
		//bad car id
		log.Printf("%s: %v for bad car_id param '%s'", tag, err, param)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pubCar, err := db.getPublicCarForID(uint(paramCarID))
	if err != nil {
		log.Printf("%s: %v for car_id param '%s'", tag, err, param)

		rspErr := http.StatusInternalServerError
		if err == errNotFound {
			rspErr = http.StatusNotFound
		}

		w.WriteHeader(rspErr)
		return
	}

	respondJSON(tag, w, pubCar)
	/*pubCarBytes, err := json.Marshal(pubCar)
	if err != nil {
		log.Printf("%s: %v for car_id param '%s'", tag, err, param)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(pubCarBytes)*/
}

//carsHandler serves path /cars
func getPublicCars(w http.ResponseWriter, r *http.Request) {
	tag := "handler.cars"
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pubCars, err := db.listPublicCars()
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respondJSON(tag, w, pubCars)
}
