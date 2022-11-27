package main

import (
	"log"
	"os"

	"github.com/svex99/bind-api/api"
	"github.com/svex99/bind-api/services/bind"
)

func main() {
	log.Println("Starting BIND API...")

	bind.Service.Init()

	router := api.SetupRouter(true)

	address := ":2020"
	if envPort := os.Getenv("PORT"); envPort != "" {
		address = ":" + envPort
	}

	if err := router.Run(address); err != nil {
		log.Fatal(err)
	}
}
