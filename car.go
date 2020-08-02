package main

import "github.com/jinzhu/gorm"

const carsTableName = "cars"

type car struct {
	gorm.Model
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
}

type pubCar struct {
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
}

type charge struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
	Per      string  `json:"per"`
}
