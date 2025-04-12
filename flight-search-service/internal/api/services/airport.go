package services

import (
	"aerona.thanhtd.com/flight-search-service/internal/api/models"
	"aerona.thanhtd.com/flight-search-service/internal/db/elasticsearch"
)

type AirportService struct {
	repo *elasticsearch.AirportRepository
}

func NewAirportService(repo *elasticsearch.AirportRepository) *AirportService {
	return &AirportService{
		repo: repo,
	}
}

func (s *AirportService) GetPopularAirports() ([]models.Airport, error) {
	return s.repo.GetPopularAirports()
}
