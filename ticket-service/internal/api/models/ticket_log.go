package models

import "time"

type TicketLog struct {
	TicketId   string     `json:"ticket_id"`
	Action     string     `json:"action"`
	ActionTime *time.Time `json:"action_time"`
	Note       *time.Time `json:"note"`
}
