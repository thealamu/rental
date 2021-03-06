package main

import "github.com/jinzhu/gorm"

const carsTableName = "cars"

type car struct {
	gorm.Model
	Name         string `json:"name"`
	Description  string `json:"description"`
	Charge       charge `json:"charge" gorm:"embedded"`
	Make         string `json:"make"`
	Transmission string `json:"transmission" gorm:"default:'manual'"`
	SeatCount    int    `json:"seatCount"`
	ImgURL       string `json:"imgUrl"`
	Color        string `json:"color"`
	AC           bool   `json:"ac"`
	Electric     bool   `json:"electric"`
	Convertible  bool   `json:"convertible"`
	IsPublic     bool   `json:"isPublic" gorm:"default:true"`
	Merchant     string `json:"merchant"`
}

//publicCar holds the minimum details any user cares about
type publicCar struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Charge       charge `json:"charge" gorm:"embedded"`
	Make         string `json:"make"`
	Transmission string `json:"transmission"`
	SeatCount    int    `json:"seatCount"`
	ImgURL       string `json:"imgUrl"`
	Color        string `json:"color"`
	AC           bool   `json:"ac"`
	Electric     bool   `json:"electric"`
	Convertible  bool   `json:"convertible"`
	Merchant     string `json:"merchant"`
}

type charge struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
	Per      string  `json:"per"`
}
