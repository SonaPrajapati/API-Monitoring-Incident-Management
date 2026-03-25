package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"api-monitoring-platform/internal/config"
	"api-monitoring-platform/internal/database"
	"api-monitoring-platform/internal/models"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			config.AppConfig.KafkaBroker,
		},
		Topic:   config.AppConfig.KafkaTopic,
		GroupID: "api-monitor-group",

		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("Kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down")
			reader.Close()
			return

		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Kafka read error:", err)
				continue
			}

			var event models.Metric
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Println("JSON decode error:", err)
				continue
			}

			collection := database.DB.Collection("consumer_metrics")
			ctxDB, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err = collection.InsertOne(ctxDB, event)
			cancel()

			if err != nil {
				log.Println("Mongo insert error:", err)
				continue
			}

			log.Println("Event stored in MongoDB")
		}
	}
}
