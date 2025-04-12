package services

import (
	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/db/elasticsearch"
)

type AirlineService struct {
	repo *elasticsearch.AirlineRepository
}

func NewAirlineService(repo *elasticsearch.AirlineRepository) *AirlineService {
	return &AirlineService{repo: repo}
}

func (s *AirlineService) GetAllAirlines() ([]models.Airline, error) {
	return s.repo.GetAllAirline()
}

func (s *AirlineService) FindByAirlineId(airlineId string) (*models.Airline, error) {
	return s.repo.FindByAirlineId(airlineId)
}

func (s *AirlineService) DeleteByAirlineId(airlineId string) error {
	return s.repo.DeleteByAirlineId(airlineId)
}
