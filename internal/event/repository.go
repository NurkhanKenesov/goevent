package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEvent(ctx context.Context, e *Event) error {
	query := `INSERT INTO events (title, description, event_date, location_type, location_data, author_id, created_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING id`
	e.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx, query,
		e.Title, e.Description, e.EventDate, e.LocationType, e.LocationData, e.AuthorID, e.CreatedAt).Scan(&e.ID)
}

func (r *Repository) GetEventByID(ctx context.Context, id int) (*Event, error) {
	e := &Event{}
	query := `SELECT id, title, description, event_date, location_type, location_data, author_id, created_at, updated_at FROM events WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventDate,
		&e.LocationType, &e.LocationData, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("event not found")
	}
	return e, err
}

func (r *Repository) UpdateEvent(ctx context.Context, e *Event) error {
	query := `
        UPDATE events SET title=$1, description=$2, event_date=$3, location_type=$4, location_data=$5, updated_at=NOW()
        WHERE id=$6
    `
	res, err := r.db.ExecContext(ctx, query,
		e.Title, e.Description, e.EventDate, e.LocationType, e.LocationData, e.ID,
	)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id int) error {
	query := `DELETE FROM events WHERE id=$1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (r *Repository) ListEvents(ctx context.Context, limit, offset int) ([]*Event, error) {
	query := fmt.Sprintf(`
        SELECT id, title, description, event_date, location_type, location_data, author_id, created_at, updated_at
        FROM events
        ORDER BY event_date ASC
        LIMIT %d OFFSET %d
    `, limit, offset)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		e := &Event{}
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate,
			&e.LocationType, &e.LocationData, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
