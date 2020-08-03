package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//carHandler serves path /cars/{car_id}
func getSinglePublicCar(w http.ResponseWriter, r *http.Request) {
	tag := "handler.car" //used to identify this function in logs
	db, err := newDatabase(defaultDbConfig)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	carIDStr := mux.Vars(r)["car_id"]
	carID, err := strconv.Atoi(carIDStr)
	if carID < 0 || err != nil {
		//bad id
		log.Printf("%s: bad car id '%s'", tag, carIDStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pubCar, err := db.getPublicCarForID(uint(carID))
	if err != nil {
		log.Printf("%s: %v for car id '%s'", tag, err, carIDStr)

		rspErr := http.StatusInternalServerError
		if err == errNotFound {
			rspErr = http.StatusNotFound
		}

		w.WriteHeader(rspErr)
		return
	}

	data, err := json.Marshal(pubCar)
	if err != nil {
		log.Printf("%s: %v for car id '%s'", tag, err, carIDStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
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

	data, err := json.Marshal(pubCars)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
