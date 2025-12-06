package event

import "time"

type Event struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	EventDate    time.Time `json:"event_date"`
	LocationType string    `json:"location_type"`
	LocationData string    `json:"location_data"`
	AuthorID     int64     `json:"author_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
