package models

import (
	"time"

	"github.com/google/uuid"
)

type City struct {
	CityId      string     `json:"city_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	CountryName string     `json:"country_name"`
	Popularity  int64      `json:"popularity"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Status      string     `json:"status"`
}

type RawCity struct {
	CityName    string `json:"city_name"`
	IATACode    string `json:"iata_code"`
	CountryISO2 string `json:"country_iso2"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Timezone    string `json:"timezone"`
	GMT         string `json:"gmt"`
	GeonameId   string `json:"geoname_id"`
}

func ParseCity(rawData RawCity) City {
	now := time.Now()
	data := City{
		CityId:      uuid.New().String(),
		Name:        rawData.CityName,
		Code:        rawData.IATACode,
		CountryName: "", //TODO
		Popularity:  0,
		CreatedAt:   &now,
		Status:      "active",
	}
	return data
}
