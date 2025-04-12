package main

import (
	"fmt"
	"os"

	"aerona.thanhtd.com/airline-integration-service/internal/api/handlers"
	"aerona.thanhtd.com/airline-integration-service/internal/api/services"
	"aerona.thanhtd.com/airline-integration-service/internal/configs"
	"aerona.thanhtd.com/airline-integration-service/internal/db/elasticsearch"
	"aerona.thanhtd.com/airline-integration-service/internal/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	// godotenv.Load(".env")

	// Init logger
	logger, err := logging.NewLogger(os.Getenv("LOG_PATH"), os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Connect to Elasticsearch
	es, err := configs.NewElasticClient()
	if err != nil {
		logger.Sugar().Fatalf("Elasticsearch error: %v", err)
	}
	logger.Info("Successfully connect to Elasticsearch")
	// Init layers
	flightRepo := elasticsearch.NewFlightRepository(es)
	flightService := services.NewFlightService(flightRepo)
	flightHandler := handlers.NewFlightHandler(flightService, logger)

	airportRepo := elasticsearch.NewAirportRepository(es)
	airportService := services.NewAirportService(airportRepo)
	airportHandler := handlers.NewAirportHandler(airportService, logger)

	airlineRepo := elasticsearch.NewAirlineRepository(es)
	airlineService := services.NewAirlineService(airlineRepo)
	airlineHandler := handlers.NewAirlineHandler(airlineService, logger)

	cityRepo := elasticsearch.NewCityRepository(es)
	cityService := services.NewCityService(cityRepo)
	cityHandler := handlers.NewCityHandler(cityService, logger)

	monitorHandler := handlers.NewMonitorHandler(logger)

	// Upload data to Elasticsearch
	// err = cityService.ImportCityData("./data/cities.json")
	// if err != nil {
	// 	logger.Fatal("Error uploading city data: %v", zap.Error(err))
	// }

	// err = airlineService.ImportAirlineData("./data/airline.json")
	// if err != nil {
	// 	logger.Fatal("Error uploading airline data: %v", zap.Error(err))
	// }

	// err = airportService.ImportAirportData("./data/airport.json")
	// if err != nil {
	// 	logger.Fatal("Error uploading airport data: %v", zap.Error(err))
	// }
	// err = flightService.ImportFlightData("./data/flights.json")
	// if err != nil {
	// 	logger.Fatal("Failed to upload flight data: %v", zap.Error(err))
	// }

	// Init router
	server := gin.Default()
	server.SetTrustedProxies(nil)
	server.POST("/api/v1/flights", flightHandler.CreateFlight)
	server.GET("/api/v1/flights", flightHandler.GetAllFlights)
	server.GET("/api/v1/flights/:flightId", flightHandler.FindByFlightId)
	// server.PUT("/api/v1/flights/:flightId", flightHandler.UpdateFlight)
	server.DELETE("/api/v1/flights/:flightId", flightHandler.DeleteByFlightId)

	server.POST("/api/v1/airports", airportHandler.CreateAirport)
	server.GET("/api/v1/airports", airportHandler.GetAllAirports)
	server.GET("/api/v1/airports/:airportId", airportHandler.FindByAirportId)
	// server.PUT("/api/v1/airports/:airportId", airportHandler.UpdateAirport)
	server.DELETE("/api/v1/airports/:airportId", airportHandler.DeleteByAirportId)

	server.POST("/api/v1/airlines", airlineHandler.CreateAirline)
	server.GET("/api/v1/airlines", airlineHandler.GetAllAirlines)
	server.GET("/api/v1/airlines/:airlineId", airlineHandler.FindByAirlineId)
	server.PUT("/api/v1/airlines/:airlineId", airlineHandler.UpdateAirline)
	// server.DELETE("/api/v1/airlines/:airlineId", airlineHandler.DeleteByAirlineId)

	server.POST("/api/v1/cities", cityHandler.CreateCity)
	server.GET("/api/v1/cities", cityHandler.GetAllCities)
	server.GET("/api/v1/cities/:cityId", cityHandler.FindByCityId)
	// server.PUT("/api/v1/cities/:cityId", cityHandler.UpdateCity)
	server.DELETE("/api/v1/cities/:cityId", cityHandler.DeleteByCityId)

	server.GET("/health", monitorHandler.HealthCheck)

	port := os.Getenv("PORT")
	logger.Info("Starting server", zap.String("port", port))
	if server.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal("Error running server: %v", zap.Error(err))
	}
}
