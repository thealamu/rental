package main

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var testDBConfig = &dbconfig{
	dialect: "sqlite3",
	dbURI:   "file::memory:?cache=shared",
}

func TestGetPublicCarForID(t *testing.T) {
	db, err := newDatabase(testDBConfig)
	if err != nil {
		t.Error(err)
	}

	//init
	var firstCar car
	firstCar.ID = 11
	var testCar car
	testCar.ID = 12

	//insert
	db.gormDB.Create(firstCar)
	db.gormDB.Create(testCar)
	//read
	savedCar, err := db.getPublicCarForID(testCar.ID)
	if err != nil {
		t.Error(err)
	}

	//validate
	if savedCar.ID != testCar.ID {
		t.Errorf("getPublicCarForID returns incorrect car, expected car with ID %d, got ID %d", testCar.ID, savedCar.ID)
	}
}

func TestSetupGDBErrorOnBadURI(t *testing.T) {
	err := setupGDB("mysql", "baduri")
	if err == nil {
		t.Error("Expected error from setupGDB for bad uri, got nil")
	}
}

func TestNewDB(t *testing.T) {
	db, err := newDatabase(testDBConfig)
	if err != nil {
		t.Error(err)
	}
	if db.gormDB == nil {
		t.Errorf("Expected returned database to have a gormDB, got nil")
		return
	}

	db.close()
}
