package main

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestNewDB(t *testing.T) {
	var c dbconfig
	c.dialect = "sqlite3"
	c.dbURI = "file::memory:?cache=shared"

	db, err := newDatabase(&c)
	if err != nil {
		t.Error(err)
	}
	if db.gormDB == nil {
		t.Errorf("Expected returned database to have a gormDB, got nil")
		return
	}

	db.close()
}
