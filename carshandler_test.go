package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSingleCarHandler(t *testing.T) {
	testRequest, err := http.NewRequest(http.MethodGet, "/cars/1", nil)
	if err != nil {
		t.Error(err)
	}

	router := mux.NewRouter()
	carsRouter := router.PathPrefix("/cars").Subrouter()
	carsRouter.HandleFunc("/{car_id:[0-9]+}", getSinglePublicCar).Methods(http.MethodGet)

	//test 500 for db error
	defaultDbConfig = &dbconfig{dialect: "mysql", dbURI: "''@this.that"}
	errRecorder := httptest.NewRecorder()

	router.ServeHTTP(errRecorder, testRequest)
	if errRecorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusInternalServerError, errRecorder.Code)
	}

	//test 200
	var c dbconfig
	c.dialect = "sqlite3"
	c.dbURI = "file::memory:?cache=shared"
	defaultDbConfig = &c

	//insert test data
	testCar := car{}
	testDb, err := newDatabase(defaultDbConfig)
	if err != nil {
		t.Error(err)
	}
	testDb.gormDB.Create(&testCar)

	rRecorder := httptest.NewRecorder()
	router.ServeHTTP(rRecorder, testRequest)
	if rRecorder.Code != http.StatusOK {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusOK, rRecorder.Code)
	}

	//test bad car id doesn't pass throught router
	badRequest, err := http.NewRequest(http.MethodGet, "/cars/-9", nil)
	if err != nil {
		t.Error(err)
	}
	errRecorder = httptest.NewRecorder()

	router.ServeHTTP(errRecorder, badRequest)
	if errRecorder.Code != http.StatusNotFound {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusNotFound, errRecorder.Code)
	}

	//test car id not found
	notFoundRequest, err := http.NewRequest(http.MethodGet, "/cars/10", nil)
	if err != nil {
		t.Error(err)
	}
	notFoundRecorder := httptest.NewRecorder()

	router.ServeHTTP(notFoundRecorder, notFoundRequest)
	if notFoundRecorder.Code != http.StatusNotFound {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusNotFound, notFoundRecorder.Code)
	}
}

func TestCarsHandler(t *testing.T) {
	testRequest, err := http.NewRequest(http.MethodGet, "/cars", nil)
	router := mux.NewRouter()
	carsRouter := router.PathPrefix("/cars").Subrouter()
	carsRouter.HandleFunc("", getPublicCars).Methods(http.MethodGet)

	//test 500 for db error
	gdb = nil
	defaultDbConfig = &dbconfig{dialect: "mysql", dbURI: "''@somepath"}

	errRecorder := httptest.NewRecorder()
	router.ServeHTTP(errRecorder, testRequest)

	if errRecorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusInternalServerError, errRecorder.Code)
	}

	//use sqlite as default db
	var c dbconfig
	c.dialect = "sqlite3"
	c.dbURI = "file::memory:?cache=shared"
	defaultDbConfig = &c

	if err != nil {
		t.Error(err)
	}
	rRecorder := httptest.NewRecorder()

	router.ServeHTTP(rRecorder, testRequest)
	if rRecorder.Code != http.StatusOK {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusOK, rRecorder.Code)
	}
}
