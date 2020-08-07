package main

import "github.com/jinzhu/gorm"

const merchantsTableName = "merchants"

type merchant struct {
	gorm.Model
	Email       string
	Name        string `gorm:"primary_key"`
	Website     string
	PublicCars  int
	PrivateCars int
}

type minimalMerchant struct {
	Name       string
	Website    string
	PublicCars int
}
