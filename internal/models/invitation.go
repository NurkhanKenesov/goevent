package models

import "time"

type Invitation struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // "pending", "accepted", "declined"
	CreatedAt time.Time `json:"created_at"`
}
