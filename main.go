package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tjeastmond/future-take-home/routes"
	"github.com/tjeastmond/future-take-home/store"
)

func main() {
	store.Connect()
	appStore := &store.Store{
		Appointments: &store.Appointments{},
	}

	router := gin.Default()
	routes.InitializeRoutes(router, appStore)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	c := router.Group("/")
	{
		c.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Hello, Future!",
			})
		})
	}

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, closing database connection...")
		store.CloseDB()
		os.Exit(0)
	}()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
