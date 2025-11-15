package repository

import (
	"goevent/internal/models"

	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	query := `
        INSERT INTO events (title, description, date, location, creator_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	return r.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.Date,
		event.Location,
		event.CreatorID,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&event.ID)
}

func (r *EventRepository) GetByID(id int) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM events WHERE id = $1`
	err := r.db.Get(&event, query, id)
	return &event, err
}

func (r *EventRepository) GetByCreatorID(creatorID int) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT * FROM events WHERE creator_id = $1 ORDER BY date DESC`
	err := r.db.Select(&events, query, creatorID)
	return events, err
}

func (r *EventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	query := `SELECT * FROM events ORDER BY date DESC`
	err := r.db.Select(&events, query)
	return events, err
}

func (r *EventRepository) Update(event *models.Event) error {
	query := `
        UPDATE events
        SET title = $1, description = $2, date = $3, location = $4, updated_at = $5
        WHERE id = $6 AND creator_id = $7
    `
	_, err := r.db.Exec(
		query,
		event.Title,
		event.Description,
		event.Date,
		event.Location,
		event.UpdatedAt,
		event.ID,
		event.CreatorID,
	)
	return err
}

func (r *EventRepository) Delete(id, creatorID int) error {
	query := `DELETE FROM events WHERE id = $1 AND creator_id = $2`
	_, err := r.db.Exec(query, id, creatorID)
	return err
}
