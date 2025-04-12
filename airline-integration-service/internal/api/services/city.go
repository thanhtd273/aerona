package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"aerona.thanhtd.com/airline-integration-service/internal/db/elasticsearch"
)

type CityService struct {
	repo *elasticsearch.CityRepository
}

func NewCityService(repo *elasticsearch.CityRepository) *CityService {
	return &CityService{repo: repo}
}

func (s *CityService) CreateCity(rawData models.RawCity) (*models.City, error) {
	city := models.ParseCity(rawData)
	return s.repo.CreateCity(city)
}

func (s *CityService) FindByCityId(cityId string) (*models.City, error) {
	return s.repo.FindByCityId(cityId)
}

func (s *CityService) GetAllCities() ([]models.City, error) {
	return s.repo.GetAllCity()
}

func (s *CityService) UpdateCity(cityId string, rawData models.RawCity) (*models.City, error) {
	return s.repo.UpdateCity(cityId, rawData)
}

func (s *CityService) DeleteByCityId(cityId string) error {
	return s.repo.DeleteByCityId(cityId)
}

func (s *CityService) ImportCityData(filepath string) error {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to read json file: %v", err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read byte value: %v", err.Error())
	}
	var cityData []models.RawCity
	err = json.Unmarshal(byteValue, &cityData)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err.Error())
	}

	for index, rawCity := range cityData {
		city := models.ParseCity(rawCity)
		s.repo.CreateCity(city)
		fmt.Printf("Successfully created city %v/%v cities\n", index, len(cityData))
	}

	return nil
}
