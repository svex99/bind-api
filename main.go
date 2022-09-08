package main

import (
	"log"

	"github.com/svex99/bind-api/api"
	"github.com/svex99/bind-api/models"
)

func main() {
	models.ConnectDatabase()

	router := api.SetupRouter()

	if err := router.Run(":2020"); err != nil {
		log.Fatal(err)
	}
}
