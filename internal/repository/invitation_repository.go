package repository

import (
	"goevent/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type InvitationRepository struct {
	db *sqlx.DB
}

func NewInvitationRepository(db *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Create(invitation *models.Invitation) error {
	query := `
        INSERT INTO invitations (event_id, invitee_id, inviter_id, status, message, sent_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `
	return r.db.QueryRow(
		query,
		invitation.EventID,
		invitation.InviteeID,
		invitation.InviterID,
		invitation.Status,
		invitation.Message,
		invitation.SentAt,
		invitation.CreatedAt,
		invitation.UpdatedAt,
	).Scan(&invitation.ID)
}

func (r *InvitationRepository) GetByID(id int) (*models.Invitation, error) {
	var invitation models.Invitation
	query := `SELECT * FROM invitations WHERE id = $1`
	err := r.db.Get(&invitation, query, id)
	return &invitation, err
}

func (r *InvitationRepository) GetByEventID(eventID int) ([]models.Invitation, error) {
	var invitations []models.Invitation
	query := `SELECT * FROM invitations WHERE event_id = $1 ORDER BY sent_at DESC`
	err := r.db.Select(&invitations, query, eventID)
	return invitations, err
}

func (r *InvitationRepository) GetByInviteeID(inviteeID int) ([]models.Invitation, error) {
	var invitations []models.Invitation
	query := `SELECT * FROM invitations WHERE invitee_id = $1 ORDER BY sent_at DESC`
	err := r.db.Select(&invitations, query, inviteeID)
	return invitations, err
}

func (r *InvitationRepository) GetByInviterID(inviterID int) ([]models.Invitation, error) {
	var invitations []models.Invitation
	query := `SELECT * FROM invitations WHERE inviter_id = $1 ORDER BY sent_at DESC`
	err := r.db.Select(&invitations, query, inviterID)
	return invitations, err
}

func (r *InvitationRepository) UpdateStatus(id int, status models.InvitationStatus, respondedAt *time.Time) error {
	query := `
        UPDATE invitations
        SET status = $1, responded_at = $2, updated_at = $3
        WHERE id = $4
    `
	_, err := r.db.Exec(query, status, respondedAt, time.Now(), id)
	return err
}

func (r *InvitationRepository) Delete(id int) error {
	query := `DELETE FROM invitations WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *InvitationRepository) Exists(eventID, inviteeID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM invitations WHERE event_id = $1 AND invitee_id = $2`
	err := r.db.Get(&count, query, eventID, inviteeID)
	return count > 0, err
}

func (r *InvitationRepository) GetByEventAndInvitee(eventID, inviteeID int) (*models.Invitation, error) {
	var invitation models.Invitation
	query := `SELECT * FROM invitations WHERE event_id = $1 AND invitee_id = $2`
	err := r.db.Get(&invitation, query, eventID, inviteeID)
	return &invitation, err
}
