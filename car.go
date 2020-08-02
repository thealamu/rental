package main

import "github.com/jinzhu/gorm"

type car struct {
	gorm.Model
	Name        string
	Description string
	Charge      charge
}

type charge struct {
	amount   float32
	currency string
	per      string
}
