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

//listPublicCars returns all cars made public by merchants
func (d database) listPublicCars() ([]pubCar, error) {
	var pubs []pubCar
	err := d.gormDB.Table(carsTableName).Find(&pubs).Error
	return pubs, err
}

//close closes the gorm db
func (d database) close() {
	d.gormDB.Close()
}

var (
	once sync.Once
	gdb  *gorm.DB
)

//newDatabase returns a new db object using the same underlying gorm db
func newDatabase() (ret database, err error) {
	ret = database{}
	once.Do(func() {
		dbURI := os.ExpandEnv("${RTL_USER}:${RTL_PASS}@/${RTL_DB}?charset=utf8mb4&parseTime=True&loc=Local")
		err = setupGDB(dbURI)
	})
	ret.gormDB = gdb
	return
}

//setupGDB does init for the gorm db
func setupGDB(dbURI string) error {
	conn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		log.Fatal("setupGormDB: ", err)
	}
	gdb = conn
	gdb.AutoMigrate(&car{})
}
