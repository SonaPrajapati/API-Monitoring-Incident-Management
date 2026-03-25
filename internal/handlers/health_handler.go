package handlers

import (
	"context"
	"net/http"
	"time"

	"api-monitoring-platform/internal/database"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

func Ready(c *gin.Context) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	err := database.DB.Client().Ping(ctx, nil)

	if err != nil {

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "READY",
	})
}
