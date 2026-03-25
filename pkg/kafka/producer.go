package kafka

import (
	"context"
	"log"
	"time"

	"api-monitoring-platform/internal/config"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitProducer() {

	Writer = &kafka.Writer{
		Addr: kafka.TCP(
			config.AppConfig.KafkaBroker,
		),
		Topic: config.AppConfig.KafkaTopic,

		Balancer: &kafka.LeastBytes{},
	}
}

func Publish(message []byte) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := Writer.WriteMessages(
		ctx,
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		log.Println("Kafka publish error:", err)
	}
}

func CloseProducer() {

	if Writer != nil {
		Writer.Close()
	}
}
