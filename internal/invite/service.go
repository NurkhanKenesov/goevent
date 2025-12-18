package invite

import (
	"context"
	"math"
	"time"

	"goevent/internal/errors"
	"goevent/internal/event"
	"goevent/internal/logging"
	"goevent/internal/middleware"
)

type InviteJob struct {
	InvitationID int64
	Email        string
	EventID      int64
	UserID       int64
	RequestID    string
	Attempts     int
}

type Service struct {
	repo       *Repository
	eventRepo  *event.Repository
	jobs       chan InviteJob
	maxRetries int
}

func NewService(ctx context.Context, repo *Repository, eventRepo *event.Repository) *Service {
	s := &Service{repo: repo, eventRepo: eventRepo, jobs: make(chan InviteJob, 100), maxRetries: 5}
	for i := 0; i < 4; i++ {
		go s.worker(ctx)
	}
	return s
}

func (s *Service) InviteUser(ctx context.Context, inviterID, inviteeID, eventID int64) error {
	e, err := s.eventRepo.GetEventByID(ctx, int(eventID))
	if err != nil {
		return errors.ErrInternal(err.Error())
	}
	if e.AuthorID != inviterID {
		return errors.ErrForbidden("only event author can invite users")
	}

	id, err := s.repo.CreateInvitation(ctx, inviteeID, eventID)
	if err != nil {
		return errors.ErrInternal(err.Error())
	}

	reqID := ""
	if v := ctx.Value("request_id"); v != nil {
		if sID, ok := v.(string); ok {
			reqID = sID
		}
	}
	job := InviteJob{InvitationID: id, EventID: eventID, UserID: inviteeID, RequestID: reqID}
	select {
	case s.jobs <- job:
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.ErrInternal("invite queue full, please retry")
	}
	return nil
}

func (s *Service) RespondInvitation(ctx context.Context, userID, invitationID int64, status Status) error {
	if status != StatusAccepted && status != StatusDeclined {
		return errors.ErrValidation("invalid status: must be 'accepted' or 'declined'")
	}
	return s.repo.UpdateInvitationStatus(ctx, invitationID, status)
}

func (s *Service) GetMyEvents(ctx context.Context, userID int64) ([]*event.Event, error) {
	return s.repo.GetMyEvents(ctx, userID)
}

func (s *Service) worker(ctx context.Context) {
	logging.Info("invite.worker.start", nil)
	for {
		select {
		case <-ctx.Done():
			logging.Info("invite.worker.stop", nil)
			return
		case job := <-s.jobs:
			s.processJob(ctx, job)
		}
	}
}

func (s *Service) processJob(ctx context.Context, job InviteJob) {
	ctxJob := context.WithValue(context.Background(), middleware.RequestIDKey, job.RequestID)
	attempt := 0
	for {
		attempt++
		fields := map[string]interface{}{
			"invitation_id": job.InvitationID,
			"event_id":      job.EventID,
			"user_id":       job.UserID,
			"attempt":       attempt,
		}
		if rid := middleware.FromContext(ctxJob); rid != "" {
			fields["request_id"] = rid
		}

		// Attempt to send invite (simulate external API call)
		logging.Info("invite.send.attempt", fields)
		
		// Simulate occasional failure (for demo purposes)
		success := attempt >= 2 // Fail first attempt, succeed on retry
		if success {
			logging.Info("invite.send.success", fields)
			return
		}

		// Retry logic with exponential backoff
		if attempt >= s.maxRetries {
			fields["reason"] = "max_retries_exceeded"
			logging.Error("invite.send.failed", fields)
			return
		}

		// Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
		backoffMs := int64(math.Pow(2, float64(attempt-1))) * 100
		fields["backoff_ms"] = backoffMs
		logging.Info("invite.send.retry", fields)

		select {
		case <-ctx.Done():
			logging.Warn("invite.send.cancelled", fields)
			return
		case <-time.After(time.Duration(backoffMs) * time.Millisecond):
			// Continue to next attempt
		}
	}
}
