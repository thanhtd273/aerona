package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aerona.thanhtd.com/ticket-service/internal/api/handlers"
	"aerona.thanhtd.com/ticket-service/internal/api/services"
	"aerona.thanhtd.com/ticket-service/internal/configs"
	"aerona.thanhtd.com/ticket-service/internal/db/mongodb"
	"aerona.thanhtd.com/ticket-service/internal/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to load environment variables: %v", err))
	// }
	// Logging
	logger, err := logging.NewLogger(os.Getenv("LOG_PATH"), os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// Context cho graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// MongoDB
	client := configs.NewMongoClient()
	defer client.Disconnect(context.Background())

	// AWS
	awsConfig := configs.NewAWSConfig(logger)
	awsClient := configs.NewS3Client(awsConfig)

	// Kafka
	consumer, err := configs.NewConsumer()
	if err != nil {
		logger.Fatal(err.Error())
	}
	producer, err := configs.NewProducer()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Init service
	ticketRepo := mongodb.NewTicketRepository(client)
	s3Service := services.NewS3Service(awsClient, os.Getenv("S3_BUCKET_NAME"), os.Getenv("AWS_REGION_ID"))
	pdfService := services.NewPDFService("templates/ticket.html")
	redisService, err := services.NewRedisClient()
	if err != nil {
		logger.Fatal(err.Error())
	}

	ticketService := services.NewTicketService(ticketRepo, s3Service, pdfService, redisService, producer, consumer, logger)
	ticketService.StartConsumer(ctx)

	ticketHandler := handlers.NewTicketHandler(ticketService, logger)
	monitorHandler := handlers.NewMonitorHandler(logger)

	server := gin.Default()
	server.SetTrustedProxies(nil)
	// server.POST("/internal/tickets", ticketHandler.CreateTicket)
	server.GET("/internal/tickets", ticketHandler.GetAllTickets)
	server.GET("/health", monitorHandler.HealthCheck)

	go func() {

		port := os.Getenv("PORT")
		if err := server.Run(fmt.Sprintf(":%s", port)); err != nil {
			logger.Fatal("Error running server: %v", zap.Error(err))
		}
		logger.Info("Starting server", zap.String("port", port))
	}()

	// Signal handling
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
