package event

import (
	"context"
	"errors"

	"goevent/internal/models"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type Service struct {
	eventRepo *Repository
	userRepo  UserRepository
}

func NewService(eventRepo *Repository, userRepo UserRepository) *Service {
	return &Service{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (s *Service) CreateEvent(
	ctx context.Context,
	e *Event,
	email string,
) error {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return ErrNotFound
	}

	e.AuthorID = user.ID
	return s.eventRepo.CreateEvent(ctx, e)
}

func (s *Service) UpdateEvent(
	ctx context.Context,
	eventID int,
	updated *Event,
	email string,
) error {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return ErrNotFound
	}

	existing, err := s.eventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		return ErrNotFound
	}

	if existing.AuthorID != user.ID {
		return ErrForbidden
	}

	updated.ID = int64(eventID)
	updated.AuthorID = user.ID

	return s.eventRepo.UpdateEvent(ctx, updated)
}

func (s *Service) DeleteEvent(
	ctx context.Context,
	eventID int,
	email string,
) error {

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return ErrNotFound
	}

	existing, err := s.eventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		return ErrNotFound
	}

	if existing.AuthorID != user.ID {
		return ErrForbidden
	}

	return s.eventRepo.DeleteEvent(ctx, eventID)

}
func (s *Service) GetEvent(ctx context.Context, id int) (*Event, error) {
	return s.eventRepo.GetEventByID(ctx, id)
}

func (s *Service) ListEvents(ctx context.Context, limit, offset int) ([]*Event, error) {
	return s.eventRepo.ListEvents(ctx, limit, offset)
}
