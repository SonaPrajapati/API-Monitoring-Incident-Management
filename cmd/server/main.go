package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-monitoring-platform/internal/config"
	"api-monitoring-platform/internal/database"
	"api-monitoring-platform/internal/handlers"
	"api-monitoring-platform/internal/middleware"
	"api-monitoring-platform/internal/worker"
	"api-monitoring-platform/pkg/circuitbreaker"
	"api-monitoring-platform/pkg/kafka"
	"api-monitoring-platform/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()
	logger.InitLogger()
	circuitbreaker.InitBreaker()
	kafka.InitProducer()

	// Connect database
	err := database.ConnectMongo()
	if err != nil {
		log.Fatal("MongoDB connection failed")
	}

	// Create root context for shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker
	go worker.StartMonitor(ctx)
	go kafka.StartConsumer(ctx)
	// Create router
	router := gin.Default()

	router.GET("/health", handlers.Health)
	router.GET("/ready", handlers.Ready)

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/apis", handlers.CreateAPI)
	protected.GET("/apis", handlers.GetAPIs)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Println("Server started on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	// Listen for OS signals
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutdown signal received")

	// Stop background worker
	cancel()

	// Create graceful shutdown timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	// Shutdown HTTP server (finish in-flight requests)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close Kafka producer after server stops
	kafka.CloseProducer()

	log.Println("Server exited gracefully")
}
