package services

import (
	"fmt"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/dto"
	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/db/elasticsearch"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"aerona.thanhtd.com/flight-search-service/internal/validator"
)

type SearchHistoryService struct {
	searchHistoryRepo *elasticsearch.SearchHistoryRepository
	airportRepo       *elasticsearch.AirportRepository
	cityRepo          *elasticsearch.CityRepository
}

func NewSearchHistoryService(searchHistoryRepo *elasticsearch.SearchHistoryRepository,
	airportRepo *elasticsearch.AirportRepository,
	cityRepo *elasticsearch.CityRepository) *SearchHistoryService {
	return &SearchHistoryService{
		searchHistoryRepo: searchHistoryRepo,
		airportRepo:       airportRepo,
		cityRepo:          cityRepo,
	}
}

func (s *SearchHistoryService) CreateOrUpdateHistory(browserId string, searchInfo dto.SearchInfo) (*models.SearchHistory, error) {
	if searchInfo.From == "" || searchInfo.To == "" {
		return nil, fmt.Errorf("departure or arrival airport is empty")
	}
	if !validator.ValidateDateString(searchInfo.DepartureDate) {
		return nil, fmt.Errorf("departure date is not a date format")
	}

	searchCode := utils.Hash(searchInfo.From, searchInfo.To, searchInfo.DepartureDate)
	isExist, err := s.searchHistoryRepo.IsHistoryExisting(searchCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence of search history, error: %s", err)
	}

	departureAirport, err := s.airportRepo.FindByAirportCode(searchInfo.From)
	if err != nil {
		return nil, fmt.Errorf("failed to find departure airport, error: %v", err)
	}
	departureCity, err := s.cityRepo.FindByCityCode(departureAirport.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to find with code=%s", departureAirport.Code)
	}

	arrivalAirport, err := s.airportRepo.FindByAirportCode(searchInfo.To)
	if err != nil {
		return nil, fmt.Errorf("failed to find arrival airport, error: %v", err)
	}
	arrivalCity, err := s.cityRepo.FindByCityCode(arrivalAirport.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to find with code=%s", arrivalAirport.Code)
	}

	departureDate, _ := time.Parse("2006-01-02", searchInfo.DepartureDate)
	now := time.Now()
	searchId, err := utils.GenerateUniqueId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Snowflake ID, error: %v", err)
	}
	searchHistory := models.SearchHistory{
		BrowserId:  browserId,
		SearchId:   searchId,
		SearchCode: searchCode,
		Departure: models.SimpleAirportInfo{
			AirportCode: departureAirport.Code,
			City:        departureCity.Name,
		},
		Arrival: models.SimpleAirportInfo{
			AirportCode: arrivalAirport.Code,
			City:        arrivalCity.Name,
		},
		DepartureDate: &departureDate,
		SearchedAt:    &now,
	}
	if isExist {
		return s.searchHistoryRepo.UpdateHistory(searchHistory)
	}

	return s.searchHistoryRepo.CreateHistory(searchHistory)
}

func (s *SearchHistoryService) GetRecentSearches() ([]models.SearchHistory, error) {
	return s.searchHistoryRepo.GetRecentSearches()
}

func (s *SearchHistoryService) DeleteBySearchId(searchId string) error {
	return s.searchHistoryRepo.DeleteBySearchId(searchId)
}
