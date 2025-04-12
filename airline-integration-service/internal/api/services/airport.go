package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"aerona.thanhtd.com/airline-integration-service/internal/db/elasticsearch"
)

type AirportService struct {
	repo *elasticsearch.AirportRepository
}

func NewAirportService(repo *elasticsearch.AirportRepository) *AirportService {
	return &AirportService{repo: repo}
}

func (s *AirportService) CreateAirport(rawData models.RawAirport) (*models.Airport, error) {
	airport := models.ParseAirport(rawData)
	return s.repo.CreateAirport(airport)
}

func (s *AirportService) FindByAirportId(airportId string) (*models.Airport, error) {
	return s.repo.FindByAirportId(airportId)
}

func (s *AirportService) GetAllAirports() ([]models.Airport, error) {
	return s.repo.GetAllAirport()
}

func (s *AirportService) UpdateAirport(airportId string, rawData models.RawAirport) (*models.Airport, error) {
	return s.repo.UpdateAirport(airportId, rawData)
}

func (s *AirportService) DeleteByAirportId(airportId string) error {
	return s.repo.DeleteByAirportId(airportId)
}

func (s *AirportService) ImportAirportData(filepath string) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to read json file: %v", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read byte value: %v", err.Error())
	}
	var airportData []models.RawAirport
	err = json.Unmarshal(byteValue, &airportData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err.Error())
	}

	for index, rawAirport := range airportData {
		airport := models.ParseAirport(rawAirport)
		s.repo.CreateAirport(airport)
		fmt.Printf("Successfully created %v/%v airports\n", index, len(airportData))
	}

	return nil
}
