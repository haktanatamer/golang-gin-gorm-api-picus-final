package main

import (
	"api-gin/package/pkg/config"
	"api-gin/package/pkg/database"
	"api-gin/package/pkg/server"
	"log"
)

// @title Gin Gorm Basket Service API
// @version 1.0
// @description Basket service api provides category-product-order.

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	if err := config.Setup(); err != nil {
		log.Fatalf("config.Setup() error: %s", err)
	}

	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup() error: %s", err)
	}

	server.Setup()
}
