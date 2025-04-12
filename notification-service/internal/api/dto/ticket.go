package dto

import "time"

type Ticket struct {
	PNR          string      `json:"pnr"`
	BookingId    int64       `json:"booking_id"`
	FlightId     string      `json:"flight_id"`
	Passengers   []Passenger `json:"passengers"`
	TicketNumber string      `json:"ticket_number"`
	Status       string      `json:"status"`
	IssuedAt     *time.Time  `json:"issuedAt"`
	CanceledAt   *time.Time  `json:"canceled_at"`
	PdfUrl       string      `json:"pdf_url"`
	QrCode       string      `json:"qr_code"`
	CreatedAt    *time.Time  `json:"created_at"`
	UpdatedAt    *time.Time  `json:"updated_at"`
}
