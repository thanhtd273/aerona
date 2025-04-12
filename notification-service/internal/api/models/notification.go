package models

import "time"

type Notification struct {
	Id        string     `json:"id"`
	BookingId string     `json:"booking_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
}
