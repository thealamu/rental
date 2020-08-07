package main

import (
	"net/http"
)

var commonEndpoints = struct {
	Car             string `json:"car_url"`
	MerchantAccount string `json:"merchant_account_url"`
	CustomerAccount string `json:"customer_account_url"`
	Merchant        string `json:"merchant_url"`
	Login           string `json:"login_url"`
}{
	Car:             "/cars/{car_id}",
	MerchantAccount: "/merchants/me",
	CustomerAccount: "/customers/me",
	Merchant:        "/merchants/{merchant}",
	Login:           "/auth/login",
}

//getCommonEndpoints serves the root path
func getCommonEndpoints(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, commonEndpoints)
}
