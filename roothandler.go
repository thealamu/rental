package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//endpoints are defined here
var endpoints = struct{}{}

//rootHandler serves the root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//return all the endpoints
	data, err := json.Marshal(endpoints)
	if err != nil {
		log.Println("rootHandler:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(data))
}
