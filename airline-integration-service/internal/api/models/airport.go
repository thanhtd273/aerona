package models

import (
	"time"

	"github.com/google/uuid"
)

type RawAirport struct {
	AirportName  string `json:"airport_name"`
	IATACode     string `json:"iata_code"`
	ICAOCode     string `json:"icao_code"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	GeoNameId    string `json:"geoname_id"`
	Timezone     string `json:"timezone"`
	GMT          string `json:"gmt"`
	PhoneNumber  string `json:"phone_number"`
	CountryName  string `json:"country_name"`
	CountryISO2  string `json:"country_iso2"`
	CityIATACode string `json:"city_iata_code"`
}

type Airport struct {
	AirportId   string     `json:"airport_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	PhoneNumber string     `json:"phone_number"`
	CityCode    string     `json:"city_code"`
	CountryName string     `json:"country_name"`
	Popularity  int64      `json:"popularity"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Status      string     `json:"status"`
}

func ParseAirport(rawData RawAirport) Airport {
	now := time.Now()
	data := Airport{
		AirportId:   uuid.New().String(),
		Name:        rawData.AirportName,
		Code:        rawData.IATACode,
		PhoneNumber: rawData.PhoneNumber,
		CityCode:    rawData.CityIATACode,
		CountryName: rawData.CountryName,
		Popularity:  0,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		Status:      "active",
	}
	return data
}
