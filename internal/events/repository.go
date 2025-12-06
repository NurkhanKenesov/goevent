package events

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateEvent добавляет новое событие в БД
func (r *Repository) CreateEvent(ctx context.Context, e *Event) error {
	query := `INSERT INTO events (title, description, event_date, location_type, location_data, creator_id, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING id`
	e.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx, query,
		e.Title, e.Description, e.EventDate, e.LocationType, e.LocationData, e.CreatorID, e.CreatedAt).Scan(&e.ID)
}
