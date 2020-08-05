package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetSingleMiniMerchant(t *testing.T) {
	//test 500 for db error
	defaultDbConfig = &dbconfig{dialect: "mysql", dbURI: "''@stuff"}
	testRequest, _ := http.NewRequest(http.MethodGet, "/merchants/somemerch", nil)

	rRecorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/merchants/{merchant}", getSingleMiniMerchant)

	router.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusInternalServerError {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusInternalServerError, rRecorder.Code)
	}

	//test not found
	defaultDbConfig = testDBConfig
	rRecorder = httptest.NewRecorder()
	router.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusNotFound {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusNotFound, rRecorder.Code)
	}

	//test 200
	gdb = nil
	testDb, _ := newDatabase(defaultDbConfig)
	testDb.gormDB.Table(merchantsTableName).Create(&minimalMerchant{Name: "somemerch"})

	rRecorder = httptest.NewRecorder()
	router.ServeHTTP(rRecorder, testRequest)

	if rRecorder.Code != http.StatusOK {
		t.Errorf("Unexpected status code, want %v, got %v", http.StatusOK, rRecorder.Code)
	}
}
