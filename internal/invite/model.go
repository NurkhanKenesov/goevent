package invite

import "time"

type Status string

const (
	StatusPending  Status = "Pending"
	StatusAccepted Status = "Accepted"
	StatusDeclined Status = "Declined"
)

type Invitation struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	EventID   int64     `json:"event_id" db:"event_id"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
