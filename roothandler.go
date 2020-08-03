package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var commonEndpoints = struct {
	Car string `json:"car_url"`
}{
	Car: "/cars/{car_id}",
}

//getCommonEndpoints serves the root path
func getCommonEndpoints(w http.ResponseWriter, r *http.Request) {
	tag := "handler.root"
	eptsBytes, err := json.Marshal(commonEndpoints)
	if err != nil {
		log.Printf("%s: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(eptsBytes)
}
