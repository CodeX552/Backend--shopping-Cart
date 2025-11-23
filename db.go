package main

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("shopping_cart.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = DB.AutoMigrate(&User{}, &Item{}, &Cart{}, &CartItem{}, &Order{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
}