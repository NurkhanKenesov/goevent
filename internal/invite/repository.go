package invite

import (
	"context"
	"goevent/internal/event"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateInvitation(ctx context.Context, userID, eventID int64) (int64, error) {
	query := `INSERT INTO invitations (user_id, event_id, status) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, userID, eventID, StatusPending).Scan(&id)
	return id, err
}

func (r *Repository) UpdateInvitationStatus(ctx context.Context, invitationID int64, status Status) error {
	query := `UPDATE invitations SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, invitationID)
	return err
}

// GetMyEvents returns events the user accepted invitations for (optimized with JOIN, no N+1)
func (r *Repository) GetMyEvents(ctx context.Context, userID int64) ([]*event.Event, error) {
	// Single JOIN query - no N+1 problem. Uses index on (user_id, status) for fast lookup.
	query := `
		SELECT DISTINCT e.id, e.title, e.description, e.event_date, e.location_type, e.location_data, e.author_id, e.created_at, e.updated_at
		FROM events e
		INNER JOIN invitations i ON e.id = i.event_id AND i.status = $2
		WHERE i.user_id = $1
		ORDER BY e.event_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, StatusAccepted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*event.Event
	for rows.Next() {
		var e event.Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.EventDate, &e.LocationType, &e.LocationData, &e.AuthorID, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, nil
}
