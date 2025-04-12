package models

import (
	"time"

	"github.com/google/uuid"
)

type RawFlightData struct {
	FlightDate   string          `json:"flight_date"`
	FlightStatus string          `json:"flight_status"`
	Departure    RawLocation     `json:"departure"`
	Arrival      RawLocation     `json:"arrival"`
	Airline      RawShortAirline `json:"airline"`
	Flight       RawFlight       `json:"flight"`
	Aircraft     RawAircraft     `json:"aircraft"`
	Live         RawLive         `json:"live"`
}

type RawLocation struct {
	Airport         string     `json:"airport"`
	Timezone        string     `json:"timezone"`
	IATA            string     `json:"iata"`
	ICAO            string     `json:"icao"`
	Terminal        string     `json:"terminal"`
	Gate            string     `json:"gate"`
	Baggage         string     `json:"baggage"`
	Delay           int        `json:"delay"`
	Scheduled       *time.Time `json:"scheduled"`
	Estimated       *time.Time `json:"estimated"`
	Actual          *time.Time `json:"actual"`
	EstimatedRunway *time.Time `json:"estimated_runway"`
	ActualRunway    *time.Time `json:"actual_runway"`
}

type RawShortAirline struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
	ICAO string `json:"icao"`
}

type RawFlight struct {
	Number     string        `json:"number"`
	IATA       string        `json:"iata"`
	ICAO       string        `json:"icao"`
	CodeShared RawCodeShared `json:"codeshared"`
}

type RawAircraft struct {
	Registration string `json:"registration"`
	IATA         string `json:"iata"`
	ICAO         string `json:"icao"`
	ICAO24       string `json:"icao24"`
}

type RawLive struct {
	Updated         *time.Time `json:"updated"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	Altitude        float64    `json:"altitude"`
	Direction       float32    `json:"direction"`
	SpeedHorizontal float32    `json:"speed_horizontal"`
	SpeedVertical   float32    `json:"speed_vertical"`
	IsGround        bool       `json:"is_ground"`
}

type RawCodeShared struct {
	AirlineNumber int    `json:"airline_number"`
	AirlineIATA   string `json:"airline_iata"`
	AirlineICAO   string `json:"airline_icao"`
	FlightNumber  string `json:"flight_number"`
	FlightIATA    string `json:"flight_iata"`
	FlightICAO    string `json:"flight_icao"`
}

type Flight struct {
	FlightId       string          `json:"flight_id"`
	FlightCode     string          `json:"flight_code"`
	Departure      AirportLocation `json:"departure"`
	Arrival        AirportLocation `json:"arrival"`
	Duration       int64           `json:"duration"`
	Price          int             `json:"price"`
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

func ParseFlight(rawData RawFlightData) Flight {
	data := Flight{}
	data.FlightId = uuid.New().String()
	data.FlightCode = rawData.Flight.IATA
	data.Departure.AirportName = rawData.Departure.Airport
	data.Departure.AirportCode = rawData.Departure.IATA
	data.Departure.City = ""
	data.Departure.Country = ""
	data.Departure.Scheduled = rawData.Departure.Scheduled
	data.Arrival.AirportName = rawData.Arrival.Airport
	data.Arrival.AirportCode = rawData.Arrival.IATA
	data.Arrival.City = ""
	data.Arrival.Country = ""
	data.Arrival.Scheduled = rawData.Arrival.Scheduled
	data.Duration = (*rawData.Arrival.Scheduled).Sub(*rawData.Departure.Scheduled).Milliseconds()
	data.Price = 1000 // TODO
	data.Currency = "VND"
	data.Airline = rawData.Airline.Name
	data.SeatsAvailable = 200 // TODO
	data.Stops = 1            // TODO
	data.StopDetails = []StopDetail{}
	data.AirplaneType = "" // TODO
	data.FlightDate = rawData.FlightDate
	now := time.Now()
	data.CreatedAt = &now
	data.UpdatedAt = &now
	data.Status = "active" // TODO
	return data
}
