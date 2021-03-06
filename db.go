package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	errNotFound = fmt.Errorf("Resource not found")
)

//configurations for the database connection
type dbconfig struct {
	dialect string
	dbURI   string
}

//Default connection configuration, can be used globally
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
//All database objects use the same underlying gorm db
func newDatabase(config *dbconfig) (ret database, err error) {
	ret = database{}
	if gdb == nil {
		err = setupGDB(config.dialect, config.dbURI)
	}
	ret.gormDB = gdb
	return
}

//getMerchantCarForID returns a car with the specified ID
func (d database) getMerchantCarForID(mchtName string, carID uint) (car, error) {
	var mchtCar car
	err := d.gormDB.
		Where("ID = ? AND merchant = ?", carID, mchtName).
		First(&mchtCar).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return mchtCar, err
}

func (d database) getNextCarID() uint {
	var lastCar car
	err := d.gormDB.Last(&lastCar).Error
	if err != nil {
		return 0
	}
	return lastCar.ID + 1
}

//update the public and private car counts for merchant
func (d database) updateMerchantCarCount(email string, isPublic bool) error {
	var colName = "private_cars"
	if isPublic {
		colName = "public_cars"
	}

	err := d.gormDB.Model(&merchant{}).
		Where("email = ?", email).
		Update(colName, gorm.Expr(colName+" + ?", 1)).Error

	return err
}

func (d database) createMerchantCar(c *car) error {
	err := d.gormDB.Create(c).Error
	return err
}

func (d database) getAuthedMerchantCars(mchtName string) []car {
	var cars []car
	d.gormDB.Find(&cars, "merchant = ?", mchtName)
	return cars
}

func (d database) searchPublicCars(query string) []publicCar {
	var pubCars []publicCar
	whereVal := fmt.Sprintf("%%%s%%", query)
	d.gormDB.Table(carsTableName).
		Where("is_public = ?", true).
		Find(&pubCars, "name LIKE ? OR make LIKE ?", whereVal, whereVal)
	return pubCars
}

func (d database) createCustomer(email string) {
	cust := customer{Email: email}
	d.gormDB.Create(&cust)
}

func (d database) mchtAccountNameExists(name string) bool {
	errMerch := d.gormDB.First(&merchant{}, "name = ?", name).Error
	return errMerch != gorm.ErrRecordNotFound
}

func (d database) userExists(email string) bool {
	errMerch := d.gormDB.First(&merchant{}, "email = ?", email).Error
	errCust := d.gormDB.First(&customer{}, "email = ?", email).Error
	return (errMerch != gorm.ErrRecordNotFound) || (errCust != gorm.ErrRecordNotFound)
}

func (d database) createMerchant(email string, name string) {
	mcht := merchant{
		Email:       email,
		Name:        name,
		PrivateCars: 0,
		PublicCars:  0,
	}
	d.gormDB.Create(&mcht)
}

func (d database) getCustomerForEmail(email string) (customer, error) {
	var cust customer
	err := d.gormDB.First(&cust, "email = ?", email).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return cust, err
}

func (d database) getMerchantForEmail(email string) (merchant, error) {
	var mcht merchant
	err := d.gormDB.First(&mcht, "email = ?", email).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return mcht, err
}

//getMerchantForName returns a merchant with the specified name
func (d database) getMiniMerchantForName(name string) (minimalMerchant, error) {
	var miniMcht minimalMerchant
	err := d.gormDB.Table(merchantsTableName).First(&miniMcht, "name = ?", name).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return miniMcht, err
}

//getPublicCarForID returns a public car with the specified ID
func (d database) getPublicCarForID(carID uint) (publicCar, error) {
	var pCar publicCar
	err := d.gormDB.Table(carsTableName).
		Where("ID = ? AND is_public = ?", carID, true).
		First(&pCar).Error
	if err == gorm.ErrRecordNotFound {
		err = errNotFound
	}
	return pCar, err
}

//listPublicCars returns all cars made public by merchants
func (d database) listPublicCars() ([]publicCar, error) {
	var pubCars []publicCar
	err := d.gormDB.Table(carsTableName).
		Where("is_public = ?", true).
		Find(&pubCars).Error
	return pubCars, err
}

//close closes the gorm db
func (d database) close() {
	d.gormDB.Close()
	gdb = nil
}

//setupGDB does init of the gorm db
func setupGDB(dialect, dbURI string) error {
	conn, err := gorm.Open(dialect, dbURI)
	if err != nil {
		return err
	}
	gdb = conn
	gdb.AutoMigrate(&car{}, &merchant{}, &customer{})

	return nil
}
