package events

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateEvent(ctx context.Context, e *Event) error {
	if strings.TrimSpace(e.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(e.Description) == "" {
		return errors.New("description is required")
	}
	if e.EventDate.IsZero() || e.EventDate.Before(time.Now()) {
		return errors.New("event_date must be in the future")
	}
	if e.LocationType != "online" && e.LocationType != "offline" {
		return errors.New("location_type must be 'online' or 'offline'")
	}
	if strings.TrimSpace(e.LocationData) == "" {
		return errors.New("location_data is required")
	}
	if e.AuthorID == 0 {
		return errors.New("author_id is required")
	}

	e.Title = strings.TrimSpace(e.Title)
	e.Description = strings.TrimSpace(e.Description)

	return s.repo.CreateEvent(ctx, e)
}
