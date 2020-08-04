package main

import "github.com/jinzhu/gorm"

type merchant struct {
	gorm.Model
	minimalMerchant
	PrivateCars int
}

type minimalMerchant struct {
	Name       string
	Website    string
	PublicCars int
}
