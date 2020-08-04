package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//respondJSON returns json data to client
func respondJSON(tag string, w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("%s: %v while encoding response", tag, err)
	}
}
