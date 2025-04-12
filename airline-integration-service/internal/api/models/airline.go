package models

import (
	"time"

	"github.com/google/uuid"
)

type RawAirline struct {
	AirlineName          string `json:"airline_name"`
	IATACode             string `json:"iata_code"`
	IATAPrefixAccounting string `json:"iata_prefix_accounting"`
	ICAOCode             string `json:"icao_code"`
	Callsign             string `json:"callsign"`
	Type                 string `json:"type"`
	Status               string `json:"status"`
	FleetSize            string `json:"fleet_size"`
	FleetAverageAge      string `json:"fleet_average_age"`
	DateFounded          string `json:"date_founded"`
	HubCode              string `json:"hub_code"`
	CountryName          string `json:"country_name"`
	CountryISO2          string `json:"country_iso2"`
}

type Airline struct {
	AirlineId   string     `json:"airline_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	CountryName string     `json:"country_name"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Status      string     `json:"status"`
}

func ParseAirline(rawData RawAirline) Airline {
	now := time.Now()
	return Airline{
		AirlineId:   uuid.New().String(),
		Name:        rawData.AirlineName,
		Code:        rawData.IATACode,
		CountryName: rawData.CountryName,
		CreatedAt:   &now,
		Status:      rawData.Status,
	}
}
