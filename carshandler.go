package main

import (
	"fmt"
	"net/http"
)

//carsHandler serves path /cars
func carsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/cars")
}
