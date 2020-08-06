package main

import "github.com/jinzhu/gorm"

type customer struct {
	gorm.Model
	Email string `json:"email" gorm:"primary_key"`
}
