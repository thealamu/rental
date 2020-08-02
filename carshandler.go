package main

import (
	"net/http"
)

//carsHandler serves path /cars
func carsHandler(w http.ResponseWriter, r *http.Request) {
	db := newDatabase()
}
