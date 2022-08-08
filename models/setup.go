package models

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	DB, err = gorm.Open(sqlite.Open("bind-api.db"), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database")
		log.Fatal("Connection Error:", err)
	} else {
		fmt.Println("Database connection is ready")
	}

	DB.AutoMigrate(&User{})
}
