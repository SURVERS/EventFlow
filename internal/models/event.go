package models

import "time"

type Event struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	EventType     string    `json:"event_type"`
	Status        string    `json:"status"`
	PublishStatus string    `json:"publish_status"` // "draft" или "published"
	CategoryID    uint      `json:"category_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEventRequest struct {
	Title         string    `json:"title" binding:"required"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	EventType     string    `json:"event_type" binding:"required"`
	Status        string    `json:"status"`
	PublishStatus string    `json:"publish_status" binding:"required"`
	CategoryID    uint      `json:"category_id" binding:"required"`
}
