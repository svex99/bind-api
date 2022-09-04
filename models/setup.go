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

	DB, err = gorm.Open(sqlite.Open("data/bind-api.db"), &gorm.Config{})
	DB.Exec("PRAGMA foreign_keys=ON")

	if err != nil {
		fmt.Println("Error connecting to database")
		log.Fatal("Connection Error:", err)
	} else {
		fmt.Println("Database connection is ready")
	}

	DB.AutoMigrate(
		&User{},
		&Domain{},
		&SOARecord{},
		&NSRecord{},
		&ARecord{},
		&MXRecord{},
		&TXTRecord{},
	)
}
