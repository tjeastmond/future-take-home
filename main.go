package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/routes"
	"github.com/tjeastmond/future-take-home/store"
)

func main() {
	ms := store.NewMemoryStore("./store/appointments.json")
	err := ms.LoadData()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
	}
	router := gin.Default()

	routes.InitializeRoutes(router, ms)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
