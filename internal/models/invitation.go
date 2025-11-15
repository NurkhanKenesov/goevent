package models

import "time"

type InvitationStatus string

const (
	InvitationPending   InvitationStatus = "pending"
	InvitationAccepted  InvitationStatus = "accepted"
	InvitationDeclined  InvitationStatus = "declined"
	InvitationCancelled InvitationStatus = "cancelled"
)

type Invitation struct {
	ID          int              `db:"id" json:"id"`
	EventID     int              `db:"event_id" json:"event_id"`
	InviteeID   int              `db:"invitee_id" json:"invitee_id"`
	InviterID   int              `db:"inviter_id" json:"inviter_id"`
	Status      InvitationStatus `db:"status" json:"status"`
	Message     string           `db:"message" json:"message"`
	SentAt      time.Time        `db:"sent_at" json:"sent_at"`
	RespondedAt *time.Time       `db:"responded_at" json:"responded_at,omitempty"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at" json:"updated_at"`
}

// Request structures
type CreateInvitationRequest struct {
	EventID   int    `json:"event_id" binding:"required"`
	InviteeID int    `json:"invitee_id" binding:"required"`
	Message   string `json:"message"`
}

type UpdateInvitationStatusRequest struct {
	Status InvitationStatus `json:"status" binding:"required,oneof=accepted declined"`
}

// Response structures
type InvitationResponse struct {
	Invitation *Invitation `json:"invitation"`
	Event      *Event      `json:"event,omitempty"`
	Inviter    *User       `json:"inviter,omitempty"`
	Invitee    *User       `json:"invitee,omitempty"`
}

type InvitationWithDetails struct {
	Invitation *Invitation `json:"invitation"`
	Event      *Event      `json:"event"`
	Inviter    *User       `json:"inviter"`
	Invitee    *User       `json:"invitee"`
}
