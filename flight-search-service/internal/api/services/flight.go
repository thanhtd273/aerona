package services

import (
	"aerona.thanhtd.com/flight-search-service/internal/api/dto"
	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/db/elasticsearch"
	"go.uber.org/zap"
)

type FlightService struct {
	logger               *zap.Logger
	flightRepo           *elasticsearch.FlightRepository
	airportRepo          *elasticsearch.AirportRepository
	cityRepo             *elasticsearch.CityRepository
	searchHistoryService *SearchHistoryService
}

func NewFlightService(logger *zap.Logger,
	flightRepo *elasticsearch.FlightRepository,
	airportRepo *elasticsearch.AirportRepository,
	cityRepo *elasticsearch.CityRepository,
	searchHistoryService *SearchHistoryService) *FlightService {
	return &FlightService{
		logger:               logger,
		flightRepo:           flightRepo,
		airportRepo:          airportRepo,
		cityRepo:             cityRepo,
		searchHistoryService: searchHistoryService,
	}
}

func (s *FlightService) SearchFlights(browserId string, searchInfo dto.SearchInfo, offset int, limit int) ([]models.Flight, error) {
	// Increase airport's popularity
	// err := s.airportRepo.IncreasePopularity(searchInfo.From)
	// if err != nil {
	// 	s.logger.Error("Failed to increase airport's popularity, error: %v", zap.Error(err))
	// }
	// err = s.airportRepo.IncreasePopularity(searchInfo.To)
	// if err != nil {
	// 	s.logger.Error("Failed to increase airport's popularity, error: %v", zap.Error(err))
	// }
	// Increase city's popularity

	// Create or update search history
	// _, err2 := s.searchHistoryService.CreateOrUpdateHistory(browserId, searchInfo)
	// if err2 != nil {
	// 	s.logger.Error("Failed to create or update search history, error: %v", zap.Error(err2))
	// }

	return s.flightRepo.SearchFlights(searchInfo, offset, limit)
}

func (s *FlightService) FilterFlights(searchInfo dto.SearchInfo) ([]models.Flight, error) {

	return s.flightRepo.FilterFlights(searchInfo)
}

func (s *FlightService) GetAllFlights() ([]models.Flight, error) {
	return s.flightRepo.GetAllFlight()
}

func (s *FlightService) FindByFlightId(flightId string) (*models.Flight, error) {
	return s.flightRepo.FindByFlightId(flightId)
}
