package main

import (
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type database struct {
	gormDB *gorm.DB
}

//configurations for the database
type dbconfig struct {
	dialect string
	dbURI   string
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

//newDatabase returns a new database object using the dialect in the config;
//it defaults to mysql. All database objects use the same underlying gorm db
func newDatabase(config *dbconfig) (ret database, err error) {
	ret = database{}
	once.Do(func() {
		if config == nil {
			//default config
			config = defaultConfig()
		}
		err = setupGDB(config.dialect, config.dbURI)
	})
	ret.gormDB = gdb
	return
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

//defaultConfig returns a mysql config
func defaultConfig() *dbconfig {
	uri := os.ExpandEnv("${RTL_USER}:${RTL_PASS}@/${RTL_DB}?charset=utf8mb4&parseTime=True&loc=Local")
	return &dbconfig{
		dbURI:   uri,
		dialect: "mysql",
	}
}
