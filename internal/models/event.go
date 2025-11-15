package models

import "time"

type Event struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
	Location    string    `db:"location" json:"location"`
	CreatorID   int       `db:"creator_id" json:"creator_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=255"`
	Description string    `json:"description" binding:"max=1000"`
	Date        time.Time `json:"date" binding:"required"`
	Location    string    `json:"location" binding:"required,max=255"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	Location    *string    `json:"location,omitempty"`
}

type EventResponse struct {
	Event *Event `json:"event"`
	User  *User  `json:"creator,omitempty"`
}
