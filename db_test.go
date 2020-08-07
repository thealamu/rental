package main

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var testDBConfig = &dbconfig{
	dialect: "sqlite3",
	dbURI:   "file::memory:?cache=shared",
}

func TestGetNextCarID(t *testing.T) {
	db, err := newDatabase(testDBConfig)
	if err != nil {
		t.Error(err)
	}

	var testCar car
	testCar.Name = "Test"
	testCar.ID = 180

	db.gormDB.Create(&testCar)

	next := db.getNextCarID()
	if next != 181 {
		t.Errorf("Expected next car ID 181, got %d", next)
	}
}

func TestGetMerchantForEmail(t *testing.T) {
	db, err := newDatabase(testDBConfig)
	if err != nil {
		t.Error(err)
	}

	var mcht merchant
	mcht.Name = "Avis Rentals"
	mcht.Email = "avis@rentals.com"

	db.gormDB.Create(&mcht)

	savedMcht, err := db.getMerchantForEmail(mcht.Email)
	if err != nil {
		t.Error(err)
	}

	if savedMcht.Email != mcht.Email {
		t.Errorf("getMerchantForName returns incorrect merchant, expected merchant with email %s, got email %s",
			mcht.Email, savedMcht.Email)
	}

	//test not found
	_, err = db.getMerchantForEmail("Some email")
	if err != errNotFound {
		t.Errorf("getMerchantForName returns %v error, expected errNotFound", err)
	}

}

func TestGetMiniMerchantForName(t *testing.T) {
	db, err := newDatabase(testDBConfig)
	if err != nil {
		t.Error(err)
	}

	//init
	var oneMerchant merchant
	oneMerchant.Name = "Anvil Car Rentals"

	//insert
	db.gormDB.Create(&oneMerchant)

	//read
	savedMcht, err := db.getMiniMerchantForName(oneMerchant.Name)
	if err != nil {
		t.Error(err)
	}

	//validate
	if savedMcht.Name != oneMerchant.Name {
		t.Errorf("getMiniMerchantForName returns incorrect merchant, expected merchant with name %s, got name %s",
			oneMerchant.Name, savedMcht.Name)
	}

	//test not found
	_, err = db.getMiniMerchantForName("Some Name")
	if err != errNotFound {
		t.Errorf("getMiniMerchantForName returns %v error, expected errNotFound", err)
	}

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
	db.gormDB.Create(&firstCar)
	db.gormDB.Create(&testCar)
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
