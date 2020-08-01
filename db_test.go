package main

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	//setup test env
	os.Setenv("RTL_DB", "rental")
	os.Setenv("RTL_USER", "rental")
	os.Setenv("RTL_PASS", "rental")

	db := newDatabase()
	if db.gormDB == nil {
		t.Errorf("Expected returned database to have a gormDB, got nil")
	}

	db.close()
}
