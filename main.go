package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/routes"
	"github.com/tjeastmond/future-take-home/store"
)

func main() {
	if err := store.LoadData(); err != nil {
		log.Printf("Warning: Failed to load data from file: %v", err)
	}

	router := gin.Default()

	routes.InitializeRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
