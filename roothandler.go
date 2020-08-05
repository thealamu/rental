package main

import (
	"net/http"
)

var commonEndpoints = struct {
	Car      string `json:"car_url"`
	Merchant string `json:"merchant_url"`
	Login    string `json:"login_url"`
}{
	Car:      "/cars/{car_id}",
	Merchant: "/merchant/{merchant}",
	Login:    "/auth/login",
}

//getCommonEndpoints serves the root path
func getCommonEndpoints(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, commonEndpoints)
}
