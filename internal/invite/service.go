package invite

import (
	"context"
	"errors"
	"goevent/internal/event"
)

type Service struct {
	repo      *Repository
	eventRepo *event.Repository
}

func NewService(repo *Repository, eventRepo *event.Repository) *Service {
	return &Service{repo: repo, eventRepo: eventRepo}
}

func (s *Service) InviteUser(ctx context.Context, inviterID, inviteeID, eventID int64) error {
	// Check if inviter is the author of the event
	e, err := s.eventRepo.GetEventByID(ctx, int(eventID))
	if err != nil {
		return err
	}
	if e.AuthorID != inviterID {
		return errors.New("only event author can invite users")
	}

	// Create invitation
	_, err = s.repo.CreateInvitation(ctx, inviteeID, eventID)
	return err
}

func (s *Service) RespondInvitation(ctx context.Context, userID, invitationID int64, status Status) error {
	// TODO: Check if userID is the invitee (optional, for security)
	// For now, assume the user is authorized via middleware

	if status != StatusAccepted && status != StatusDeclined {
		return errors.New("invalid status")
	}

	return s.repo.UpdateInvitationStatus(ctx, invitationID, status)
}

func (s *Service) GetMyEvents(ctx context.Context, userID int64) ([]*event.Event, error) {
	return s.repo.GetMyEvents(ctx, userID)
}
