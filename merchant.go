package main

import "github.com/jinzhu/gorm"

type merchant struct {
	gorm.Model
	Name        string `gorm:"PRIMARY_KEY"`
	Website     string
	PublicCars  string
	PrivateCars int
}

type minimalMerchant struct {
	Name       string
	Website    string
	PublicCars int
}
