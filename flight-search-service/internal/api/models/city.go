package models

import "time"

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
