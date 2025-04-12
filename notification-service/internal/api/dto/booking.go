package dto

import "time"

type Booking struct {
	BookingId       string      `json:"booking_id"`
	PNR             string      `json:"pnr"`
	Flight          Flight      `json:"flight"`
	NumOfPassengers int         `json:"num_of_passengers"`
	UserId          string      `json:"user_id"`
	Contact         Contact     `json:"contact"`
	Passengers      []Passenger `json:"passengers"`
	TotalPrice      float32     `json:"total_price"`
	Currency        string      `json:"currency"`
	TicketUrl       string      `json:"ticket_url"`
	Status          string      `json:"status"`
}

type Contact struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Passenger struct {
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	DayOfBirth     *time.Time `json:"day_of_birth"`
	Nationality    string     `json:"nationality"`
	PassportNumber string     `json:"passport_number"`
}
