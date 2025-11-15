package service

import (
	"errors"
	"goevent/internal/models"
	"goevent/internal/repository"
	"time"
)

type InvitationService struct {
	invitationRepo *repository.InvitationRepository
	eventRepo      *repository.EventRepository
	userRepo       *repository.UserRepository
}

func NewInvitationService(
	invitationRepo *repository.InvitationRepository,
	eventRepo *repository.EventRepository,
	userRepo *repository.UserRepository,
) *InvitationService {
	return &InvitationService{
		invitationRepo: invitationRepo,
		eventRepo:      eventRepo,
		userRepo:       userRepo,
	}
}

func (s *InvitationService) Create(req *models.CreateInvitationRequest, inviterID int) (*models.Invitation, error) {
	// Check if event exists
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	// Check if inviter is the event creator
	if event.CreatorID != inviterID {
		return nil, errors.New("only event creator can send invitations")
	}

	// Check if invitee exists
	_, err = s.userRepo.GetUserByID(req.InviteeID)
	if err != nil {
		return nil, errors.New("invitee not found")
	}

	// Check if invitation already exists
	exists, err := s.invitationRepo.Exists(req.EventID, req.InviteeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("invitation already sent to this user")
	}

	// Don't allow inviting yourself
	if req.InviteeID == inviterID {
		return nil, errors.New("cannot invite yourself")
	}

	invitation := &models.Invitation{
		EventID:   req.EventID,
		InviteeID: req.InviteeID,
		InviterID: inviterID,
		Status:    models.InvitationPending,
		Message:   req.Message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.invitationRepo.Create(invitation)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *InvitationService) GetByID(id int) (*models.Invitation, error) {
	return s.invitationRepo.GetByID(id)
}

func (s *InvitationService) GetByEventID(eventID int) ([]models.Invitation, error) {
	return s.invitationRepo.GetByEventID(eventID)
}

func (s *InvitationService) GetMyInvitations(userID int) ([]models.Invitation, error) {
	return s.invitationRepo.GetByInviteeID(userID)
}

func (s *InvitationService) GetSentInvitations(userID int) ([]models.Invitation, error) {
	return s.invitationRepo.GetByInviterID(userID)
}

func (s *InvitationService) RespondToInvitation(id int, userID int, req *models.UpdateInvitationStatusRequest) (*models.Invitation, error) {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("invitation not found")
	}

	// Check if user is the invitee
	if invitation.InviteeID != userID {
		return nil, errors.New("access denied: you can only respond to your own invitations")
	}

	// Check if invitation is still pending
	if invitation.Status != models.InvitationPending {
		return nil, errors.New("invitation has already been responded to")
	}

	now := time.Now()
	invitation.Status = req.Status
	invitation.RespondedAt = &now
	invitation.UpdatedAt = now

	err = s.invitationRepo.UpdateStatus(id, req.Status, &now)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *InvitationService) CancelInvitation(id int, userID int) error {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return errors.New("invitation not found")
	}

	// Check if user is the inviter
	if invitation.InviterID != userID {
		return errors.New("access denied: you can only cancel your own invitations")
	}

	// Check if invitation is still pending
	if invitation.Status != models.InvitationPending {
		return errors.New("cannot cancel invitation that has already been responded to")
	}

	return s.invitationRepo.Delete(id)
}

func (s *InvitationService) GetInvitationWithDetails(id int) (*models.InvitationWithDetails, error) {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("invitation not found")
	}

	event, err := s.eventRepo.GetByID(invitation.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	inviter, err := s.userRepo.GetUserByID(invitation.InviterID)
	if err != nil {
		return nil, errors.New("inviter not found")
	}

	invitee, err := s.userRepo.GetUserByID(invitation.InviteeID)
	if err != nil {
		return nil, errors.New("invitee not found")
	}

	return &models.InvitationWithDetails{
		Invitation: invitation,
		Event:      event,
		Inviter:    inviter,
		Invitee:    invitee,
	}, nil
}
