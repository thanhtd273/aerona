package dto

type SearchInfo struct {
	From            string         `json:"from"`
	To              string         `json:"to"`
	DepartureDate   string         `json:"departure_date"`
	NumOfPassengers int            `json:"num_of_passengers"`
	Filters         FilterCriteria `json:"filters"`
}

type HourRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type FilterCriteria struct {
	AirlineCode    []string    `json:"airline_code"`
	DepartureRange []HourRange `json:"departure_range"`
	ArrivalRange   []HourRange `json:"arrival_range"`
}
