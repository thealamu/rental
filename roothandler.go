package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//endpoints are defined here
var endpoints = struct{}{}

//rootHandler serves the root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tag := "handler.root"
	//return all the endpoints
	data, err := json.Marshal(endpoints)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
