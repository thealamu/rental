package main

import (
	"net/http"
)

var commonEndpoints = struct {
	Car             string `json:"car_url"`
	CarsSearch      string `json:"cars_search_url"`
	MerchantAccount string `json:"merchant_account_url"`
	Merchant        string `json:"merchant_url"`
	Login           string `json:"login_url"`
}{
	Car:             "/cars/{car_id}",
	CarsSearch:      "/cars/search",
	MerchantAccount: "/merchants/me",
	Merchant:        "/merchants/{merchant}",
	Login:           "/auth/login",
}

//getCommonEndpoints serves the root path
func getCommonEndpoints(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, commonEndpoints)
}
