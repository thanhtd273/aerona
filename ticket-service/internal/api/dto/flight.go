package dto

import (
	"time"
)

type Flight struct {
	FlightId       string          `json:"flight_id"`
	FlightCode     string          `json:"flight_code"`
	Departure      AirportLocation `json:"departure"`
	Arrival        AirportLocation `json:"arrival"`
	Duration       int64           `json:"duration"`
	Price          float32         `json:"price"`
	Currency       string          `json:"currency"`
	Airline        string          `json:"airline"`
	SeatsAvailable int             `json:"seats_available"`
	Stops          int             `json:"stops"`
	StopDetails    []StopDetail    `json:"stop_details"`
	AirplaneType   string          `json:"airplane_type"`
	FlightDate     string          `json:"flight_date"`
	CreatedAt      *time.Time      `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	Status         string          `json:"status"`
}

type AirportLocation struct {
	AirportName string     `json:"airport_name"`
	AirportCode string     `json:"airport_code"`
	City        string     `json:"city"`
	Country     string     `json:"country"`
	Scheduled   *time.Time `json:"scheduled"`
}

type StopDetail struct {
	StopId        int        `json:"stop_id"`
	AirportCode   string     `json:"airport_code"`
	DepartureTime *time.Time `json:"departure_time"`
	ArrivalTime   *time.Time `json:"arrival_time"`
}
