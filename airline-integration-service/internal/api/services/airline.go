package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"aerona.thanhtd.com/airline-integration-service/internal/db/elasticsearch"
)

type AirlineService struct {
	repo *elasticsearch.AirlineRepository
}

func NewAirlineService(repo *elasticsearch.AirlineRepository) *AirlineService {
	return &AirlineService{repo: repo}
}

func (s *AirlineService) CreateAirline(rawData models.RawAirline) (*models.Airline, error) {
	airline := models.ParseAirline(rawData)
	return s.repo.CreateAirline(airline)
}

func (s *AirlineService) GetAllAirlines() ([]models.Airline, error) {
	return s.repo.GetAllAirline()
}

func (s *AirlineService) FindByAirlineId(airlineId string) (*models.Airline, error) {
	return s.repo.FindByAirlineId(airlineId)
}

func (s *AirlineService) UpdateByAirlineId(airlineId string, updatedAirline models.RawAirline) (*models.Airline, error) {
	return s.repo.UpdateAirline(airlineId, updatedAirline)
}

func (s *AirlineService) DeleteByAirlineId(airlineId string) error {
	return s.repo.DeleteByAirlineId(airlineId)
}

func (s *AirlineService) ImportAirlineData(filepath string) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to read json file: %v", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read byte value: %v", err.Error())
	}
	var airlineData []models.RawAirline
	err = json.Unmarshal(byteValue, &airlineData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err.Error())
	}

	for index, rawAirline := range airlineData {
		airline := models.ParseAirline(rawAirline)
		s.repo.CreateAirline(airline)
		fmt.Printf("Successfully created airline %v/%v airlines\n", index, len(airlineData))
	}

	return nil
}
