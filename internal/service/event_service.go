package service

import (
	"errors"
	"goevent/internal/models"
	"goevent/internal/repository"
	"time"
)

type EventService struct {
	eventRepo *repository.EventRepository
}

func NewEventService(eventRepo *repository.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

func (s *EventService) Create(req *models.CreateEventRequest, creatorID int) (*models.Event, error) {
	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		Location:    req.Location,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.eventRepo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetByID(id int) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}
	return event, nil
}

func (s *EventService) GetByCreatorID(creatorID int) ([]models.Event, error) {
	return s.eventRepo.GetByCreatorID(creatorID)
}

func (s *EventService) GetAll() ([]models.Event, error) {
	return s.eventRepo.GetAll()
}

func (s *EventService) Update(id int, req *models.UpdateEventRequest, creatorID int) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.CreatorID != creatorID {
		return nil, errors.New("access denied: you can only update your own events")
	}

	// Update only provided fields
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Date != nil {
		event.Date = *req.Date
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	event.UpdatedAt = time.Now()

	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) Delete(id, creatorID int) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return errors.New("event not found")
	}

	if event.CreatorID != creatorID {
		return errors.New("access denied: you can only delete your own events")
	}

	return s.eventRepo.Delete(id, creatorID)
}
