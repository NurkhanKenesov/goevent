package events

import (
    "context"
    "database/sql"
    "time"
)


type Repository struct {
    DB *sql.DB
}


func NewRepository(db *sql.DB) *Repository {
    return &Repository{DB: db}
}


func (r *Repository) CreateEvent(ctx context.Context, e *Event) error {
    query := `
    INSERT INTO events
    (title, description, event_date, location_type, location_data, creator_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id
    `


    err := r.DB.QueryRowContext(
        ctx,
        query,
        e.Title,
        e.Description,
        e.EventDate,
        e.LocationType,
        e.LocationData,
        e.AuthorID,
        time.Now(), 
        time.Now(), 
    ).Scan(&e.ID) 

    return err 
}
