package main

import "github.com/jinzhu/gorm"

type car struct {
	gorm.Model
	Name         string
	Description  string
	Charge       charge
	Make         string
	Transmission string
	SeatCount    int
	ImgURL       string
	Color        string
	AC           bool
	Electric     bool
	Convertible  bool
}

type pubCar struct {
	ID           uint
	Name         string
	Description  string
	Charge       charge
	Make         string
	Transmission string
	SeatCount    int
	ImgURL       string
	Color        string
	AC           bool
	Electric     bool
	Convertible  bool
}

type charge struct {
	amount   float32
	currency string
	per      string
}
