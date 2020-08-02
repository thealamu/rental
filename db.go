package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//errors
var (
	errNotFound = fmt.Errorf("Resource not found")
)

//configurations for the database
type dbconfig struct {
	dialect string
	dbURI   string
}

var defaultDbConfig = &dbconfig{
	dialect: "mysql",
	dbURI:   os.ExpandEnv("${RTL_USER}:${RTL_PASS}@/${RTL_DB}?charset=utf8mb4&parseTime=True&loc=Local"),
}

var (
	gdb *gorm.DB
)

type database struct {
	gormDB *gorm.DB
}

//newDatabase returns a new database object using the dialect in the config.
// All database objects use the same underlying gorm db
func newDatabase(config *dbconfig) (ret database, err error) {
	ret = database{}
	if gdb == nil {
		err = setupGDB(config.dialect, config.dbURI)
	}
	ret.gormDB = gdb
	return
}

//getPublicCarForID returns a public car with the specified ID
func (d database) getPublicCarForID(carID uint) (pubCar, error) {
	var pub pubCar
	err := d.gormDB.Table(carsTableName).First(&pub, carID).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return pub, err
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

//setupGDB does init of the gorm db
func setupGDB(dialect, dbURI string) error {
	conn, err := gorm.Open(dialect, dbURI)
	if err != nil {
		return err
	}
	gdb = conn
	gdb.AutoMigrate(&car{})

	return nil
}
