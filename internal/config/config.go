package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	DBName        string
	JWTSecret     string
	CheckInterval int
	KafkaBroker   string
	KafkaTopic    string
}

var AppConfig Config

func LoadConfig() {

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using system env")
	}

	interval, _ := strconv.Atoi(
		getEnv("CHECK_INTERVAL", "30"),
	)

	AppConfig = Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:        getEnv("DB_NAME", "api_monitoring"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		CheckInterval: interval,

		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "api-status"),
	}
}

func getEnv(key, defaultVal string) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultVal
	}

	return value
}
