package models

import "time"

type SearchHistory struct {
	SearchId      string `json:"search_id"`
	BrowserId     string `json:"browser_id"`
	SearchCode    string `json:"search_code"`
	Departure     SimpleAirportInfo
	Arrival       SimpleAirportInfo
	DepartureDate *time.Time `json:"departure_date"`
	SearchedAt    *time.Time `json:"created_at"`
}

type SimpleAirportInfo struct {
	AirportCode string `json:"airport_code"`
	City        string `json:"city"`
}
