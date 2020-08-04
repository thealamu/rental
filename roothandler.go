package main

import (
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
	respondJSON(tag, w, commonEndpoints)
}
