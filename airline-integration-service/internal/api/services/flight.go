package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"aerona.thanhtd.com/airline-integration-service/internal/db/elasticsearch"
)

type FlightService struct {
	repo *elasticsearch.FlightRepository
}

func NewFlightService(repo *elasticsearch.FlightRepository) *FlightService {
	return &FlightService{repo: repo}
}

func (s *FlightService) CreateFlight(rawData models.RawFlightData) (*models.Flight, error) {
	flight := models.ParseFlight(rawData)
	return s.repo.CreateFlight(flight)
}

func (s *FlightService) GetAllFlights() ([]models.Flight, error) {
	return s.repo.GetAllFlight()
}

func (s *FlightService) FindByFlightId(flightId string) (*models.Flight, error) {
	return s.repo.FindByFlightId(flightId)
}

func (s *FlightService) UpdateByFlightId(flightId string, updatedFlight models.RawFlightData) (*models.Flight, error) {
	return s.repo.UpdateFlight(flightId, updatedFlight)
}

func (s *FlightService) DeleteByFlightId(flightId string) error {
	return s.repo.DeleteByFlightId(flightId)
}

func (s *FlightService) ImportFlightData(filepath string) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to read json file: %v", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read byte value: %v", err.Error())
	}
	var flightData []models.RawFlightData
	err = json.Unmarshal(byteValue, &flightData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err.Error())
	}

	for index, rawFlight := range flightData {
		flight := models.ParseFlight(rawFlight)
		s.repo.CreateFlight(flight)
		fmt.Printf("Successfully created flight %v/%v flights\n", index, len(flightData))
	}

	return nil
}
