package main

import "github.com/jinzhu/gorm"

const merchantsTableName = "merchants"

type merchant struct {
	gorm.Model
	Email       string `gorm:"primary_key"`
	Name        string `gorm:"unique"`
	Website     string
	PublicCars  int
	PrivateCars int
}

type minimalMerchant struct {
	Name       string
	Website    string
	PublicCars int
}
