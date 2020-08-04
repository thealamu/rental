package main

import "github.com/jinzhu/gorm"

const merchantsTableName = "merchants"

type merchant struct {
	gorm.Model
	Name        string `gorm:"PRIMARY_KEY"`
	Website     string
	PublicCars  int
	PrivateCars int
}

type minimalMerchant struct {
	Name       string `gorm:"primary_key"`
	Website    string
	PublicCars int
}
