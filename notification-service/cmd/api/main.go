package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aerona.thanhtd.com/notification-service/internal/api/handlers"
	"aerona.thanhtd.com/notification-service/internal/api/services"
	"aerona.thanhtd.com/notification-service/internal/configs"
	mongodb "aerona.thanhtd.com/notification-service/internal/db/mongo"
	"aerona.thanhtd.com/notification-service/internal/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to load environment variable: %v", err))
	// }

	// Configure logging
	logger, err := logging.NewLogger(os.Getenv("LOG_PATH"), os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// Context cho graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Establish MongoDB connection
	mongoClient := configs.NewMongoClient()
	defer mongoClient.Disconnect(context.Background())

	// Connect to Redis
	redisService, err := services.NewRedisClient()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Create Kafka consumer and producer
	consumer, err := configs.NewConsumer()
	if err != nil {
		logger.Fatal(err.Error())
	}
	producer, err := configs.NewProducer()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Configure email
	emailConfig, err := configs.LoadEmailConfig()
	if err != nil {
		logger.Fatal("Failed to load email configuration", zap.Error(err))
	}
	emailService := services.NewEmailService(emailConfig)

	notificationRepo := mongodb.NewNotificationRepository(mongoClient)
	notificationService := services.NewNotificationService(notificationRepo, emailService, redisService, producer, consumer, logger)

	// Kafka listeners
	notificationService.StartConsumer(ctx)

	// Expose REST API
	notificationHandler := handlers.NewNotificationHandler(logger, notificationService)
	monitorHandler := handlers.NewMonitorHandler(logger)

	server := gin.Default()
	server.SetTrustedProxies(nil)
	server.GET("/internal/notifications/:notificationId", notificationHandler.FindById)
	server.GET("/health", monitorHandler.HealthCheck)

	go func() {
		port := os.Getenv("PORT")
		if err := server.Run(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatal("Error running server: %v", zap.Error(err))
		}
		logger.Info("Starting server", zap.String("port", port))
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down ticket-service")
	cancel()

	time.Sleep(2 * time.Second)
	consumer.Close()
	producer.Flush(5000)
	producer.Close()
}
