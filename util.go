package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Code int
	Msg  string
}

func respondError(tag string, w http.ResponseWriter, errCode int, errMsg string, statusCode int) {
	log.Printf("%s: %v", tag, errMsg)
	errRsp := &errorResponse{
		Code: errCode,
		Msg:  errMsg,
	}
	w.WriteHeader(statusCode)
	respondJSON(w, errRsp)
}

//respondJSON returns json data to client
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("respondJSON: %v", err)
	}
}
