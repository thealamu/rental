package main

import (
	"os"
	"testing"
)

//do test setup
func TestMain(m *testing.M) {
	var c dbconfig
	c.dialect = "sqlite3"
	c.dbURI = "file::memory:?cache=shared"
	defaultDbConfig = &c

	code := m.Run()
	os.Exit(code)
}
