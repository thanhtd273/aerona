package main

import (
	"fmt"
	"os"

	"aerona.thanhtd.com/flight-search-service/internal/api/handlers"
	"aerona.thanhtd.com/flight-search-service/internal/api/services"
	"aerona.thanhtd.com/flight-search-service/internal/db/elasticsearch"
	"aerona.thanhtd.com/flight-search-service/internal/logging"

	goes "github.com/elastic/go-elasticsearch/v7"
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

	esCfg := goes.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
		// APIKey:    os.Getenv("ELASTICSEARCH_API_KEY"),
	}
	es, err := goes.NewClient(esCfg)
	if err != nil {
		logger.Fatal("Error creating Elasticsearch client: %v", zap.Error(err))
	}

	// Init layers
	flightRepo := elasticsearch.NewFlightRepository(es)
	airportRepo := elasticsearch.NewAirportRepository(es)
	cityRepo := elasticsearch.NewCityRepository(es)
	searchHistoryRepo := elasticsearch.NewSearchHistoryRepository(es)

	airportService := services.NewAirportService(airportRepo)
	searchHistoryService := services.NewSearchHistoryService(searchHistoryRepo, airportRepo, cityRepo)
	flightService := services.NewFlightService(logger, flightRepo, airportRepo, cityRepo, searchHistoryService)

	flightHandler := handlers.NewFlightHandler(flightService, logger)
	airportHandler := handlers.NewAirportHandler(airportService, logger)
	searchHistoryHandler := handlers.NewSearchHistoryHandler(searchHistoryService, logger)
	monitorHandler := handlers.NewMonitorHandler(logger)

	// if err != nil {
	// 	logger.Fatal("Error importing flight data: %v", zap.Error(err))
	// }

	// Init router
	server := gin.Default()
	server.SetTrustedProxies(nil)
	server.GET("/api/v1/flights", flightHandler.GetAllFlights)
	server.GET("/api/v1/flights/:flightId", flightHandler.FindByFlightId)
	server.GET("/api/v1/flights/search", flightHandler.SearchFlights)
	server.GET("/api/v1/flights/filter", flightHandler.FilterFlights)

	server.GET("/api/v1/airports/popular", airportHandler.GetPopularAirports)

	server.GET("/api/vi/search-histories/recent", searchHistoryHandler.GetRecentSearches)
	server.DELETE("/api/v1/search-histories/:searchId", searchHistoryHandler.DeleteBySearchId)

	server.GET("/health", monitorHandler.HealthCheck)
	port := os.Getenv("PORT")
	logger.Info("Starting server", zap.String("port", port))
	if server.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal("Error running server: %v", zap.Error(err))
	}
}
