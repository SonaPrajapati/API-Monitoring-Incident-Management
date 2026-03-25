package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api-monitoring-platform/internal/config"
	"api-monitoring-platform/internal/database"
	"api-monitoring-platform/internal/models"
	"api-monitoring-platform/pkg/kafka"
	"api-monitoring-platform/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
)

func StartMonitor(ctx context.Context) {

	ticker := time.NewTicker(
		time.Duration(config.AppConfig.CheckInterval) *
			time.Second,
	)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			fmt.Println("Worker shutting down...")
			return

		case <-ticker.C:
			checkAPIs()

		}
	}
}

func checkAPIs() {

	collection := database.DB.Collection("apis")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		fmt.Println("Error fetching APIs:", err)
		return
	}

	var apis []models.API

	err = cursor.All(ctx, &apis)

	if err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	for _, api := range apis {

		go monitorAPI(api)
	}
}

func monitorAPI(api models.API) {

	start := time.Now()

	resp, err := http.Get(api.URL)

	duration := time.Since(start)

	status := 0

	if err != nil {

		fmt.Println("API DOWN:", api.Name)

	} else {

		status = resp.StatusCode

		defer resp.Body.Close()

		logger.Log.WithFields(map[string]interface{}{
			"api":        api.Name,
			"status":     status,
			"latency_ms": duration.Milliseconds(),
		}).Info("API checked")
	}

	// Store metric in MongoDB

	collection := database.DB.Collection("metrics")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	metric := models.Metric{
		APIName:   api.Name,
		Status:    status,
		Latency:   duration.Milliseconds(),
		Timestamp: time.Now(),
	}

	_, err = collection.InsertOne(ctx, metric)

	if err != nil {
		logger.Log.WithError(err).Error("Metric insert failed")
	}

	// Publish event to Kafka

	event := map[string]interface{}{
		"api_name":  api.Name,
		"url":       api.URL,
		"status":    status,
		"latency":   duration.Milliseconds(),
		"timestamp": time.Now(),
	}

	data, err := json.Marshal(event)

	if err != nil {
		fmt.Println("Kafka marshal error:", err)
		return
	}

	kafka.Publish(data)
}
