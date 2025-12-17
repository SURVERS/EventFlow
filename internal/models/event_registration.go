package models

import "time"

type EventRegistration struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EventID       uint      `gorm:"uniqueIndex:idx_participant_event" json:"event_id"`
	ParticipantID uint      `gorm:"uniqueIndex:idx_participant_event" json:"participant_id"`
	RegisteredAt  time.Time `json:"registered_at"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateEventRegistrationRequest struct {
	EventID       uint   `json:"event_id" binding:"required"`
	ParticipantID uint   `json:"participant_id" binding:"required"`
	Status        string `json:"status" binding:"required"`
}
