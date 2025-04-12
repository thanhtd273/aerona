package models

import (
	"time"
)

type Airline struct {
	AirlineId   string     `json:"airline_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	CountryName string     `json:"country_name"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Status      string     `json:"status"`
}
