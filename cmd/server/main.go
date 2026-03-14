package main

import (
	"log"

	"api-monitoring-platform/internal/database"
	"api-monitoring-platform/internal/handlers"

	"github.com/gin-gonic/gin"

	"api-monitoring-platform/internal/middleware"
)

func main() {

	err := database.ConnectMongo()
	if err != nil {
		log.Fatal("MongoDB connection failed")
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "running"})
	})

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.Run(":8080")

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/profile", handlers.Profile)
}
