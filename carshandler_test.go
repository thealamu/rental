package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCarsHandler(t *testing.T) {
	//use sqlite as default db
	var c dbconfig
	c.dialect = "sqlite3"
	c.dbURI = "file::memory:?cache=shared"
	defaultDbConfig = &c

	testRequest, err := http.NewRequest(http.MethodGet, "/cars", nil)
	if err != nil {
		t.Error(err)
	}
	rRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	carsRouter := router.PathPrefix("/cars").Subrouter()
	carsRouter.HandleFunc("", carsHandler).Methods(http.MethodGet)

	router.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusOK {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusOK, rRecorder.Code)
	}
}
