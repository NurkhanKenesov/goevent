package invite

import (
	"context"
	"errors"
	"goevent/internal/event"
	"testing"
)

type mockInviteRepo struct{}

func (m *mockInviteRepo) CreateInvitation(ctx context.Context, inviteeID, eventID int64) (int64, error) {
	return 1, nil
}

func (m *mockInviteRepo) UpdateInvitationStatus(ctx context.Context, invitationID int64, status Status) error {
	return nil
}

func (m *mockInviteRepo) GetMyEvents(ctx context.Context, userID int64) ([]*event.Event, error) {
	return []*event.Event{}, nil
}

type mockEventRepo struct{}

func (m *mockEventRepo) GetEventByID(ctx context.Context, id int) (*event.Event, error) {
	if id == 999 {
		return nil, errors.New("event not found")
	}
	return &event.Event{ID: int64(id), AuthorID: 1}, nil
}

func TestInviteService(t *testing.T) {
	svc := NewService(&mockInviteRepo{}, &mockEventRepo{})

	t.Run("Success invite", func(t *testing.T) {
		err := svc.InviteUser(context.Background(), 1, 10, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Only author can invite", func(t *testing.T) {
		err := svc.InviteUser(context.Background(), 2, 10, 1)
		if err == nil || err.Error() != "only event author can invite users" {
			t.Errorf("expected authorization error, got %v", err)
		}
	})

	t.Run("Invalid status response", func(t *testing.T) {
		err := svc.RespondInvitation(context.Background(), 1, 1, "invalid_status")
		if err == nil || err.Error() != "invalid status" {
			t.Errorf("expected invalid status error, got %v", err)
		}
	})
}
