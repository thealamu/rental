package main

import (
	"log"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type database struct {
	gormDB *gorm.DB
}

func (d database) close() {
	d.gormDB.Close()
}

var (
	once sync.Once
	gdb  *gorm.DB
)

func newDatabase() (ret database) {
	ret = database{}
	once.Do(setupGDB)
	ret.gormDB = gdb
	return
}

func setupGDB() {
	dbURI := os.ExpandEnv("${rtl_user}:${rtl_password}@/${rtl_name}?charset=utf8mb4&parseTime=True&loc=Local")
	conn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Fatal("setupGormDB: ", err)
	}
	gdb = conn
}
