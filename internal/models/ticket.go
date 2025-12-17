package models

import (
	"time"
)

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EventID       uint      `json:"event_id"`
	ParticipantID uint      `json:"participant_id"`
	TicketType    string    `json:"ticket_type"` // "free" или "paid"
	Status        string    `json:"status"`      // "active" или "canceled"
	QRCode        string    `gorm:"unique" json:"qr_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateTicketRequest struct {
	EventID       uint   `json:"event_id" binding:"required"`
	ParticipantID uint   `json:"participant_id" binding:"required"`
	TicketType    string `json:"ticket_type" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

type UpdateTicketRequest struct {
	TicketType string `json:"ticket_type" binding:"required"`
	Status     string `json:"status" binding:"required"`
}
