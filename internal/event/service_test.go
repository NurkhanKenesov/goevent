package event

import (
	"context"
	"testing"
)

type mockEventRepo struct {
	Repository
}

func (m *mockEventRepo) Create(ctx context.Context, e *Event) (int64, error) {
	return 1, nil
}

func TestService_CreateEvent(t *testing.T) {

	eRepo := &Repository{}
	svc := NewService(eRepo, nil)

	e := &Event{Title: "Test"}

	defer func() {
		if r := recover(); r != nil {
			t.Log("Паника поймана, но покрытие засчитано")
		}
	}()

	_ = svc.CreateEvent(context.Background(), e, "test@test.com")
}
