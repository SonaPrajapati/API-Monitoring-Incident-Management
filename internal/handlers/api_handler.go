package handlers

import (
	"context"
	"net/http"
	"time"

	"api-monitoring-platform/internal/database"
	"api-monitoring-platform/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateAPI(c *gin.Context) {

	var api models.API

	if err := c.BindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.DB.Collection("apis")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, api)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API created",
	})
}

func GetAPIs(c *gin.Context) {

	collection := database.DB.Collection("apis")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch APIs"})
		return
	}

	var apis []models.API

	if err := cursor.All(ctx, &apis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apis)
}
