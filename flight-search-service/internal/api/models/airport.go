package models

import (
	"time"
)

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
