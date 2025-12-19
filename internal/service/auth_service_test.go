package service

import (
	"context"
	"errors"
	"goevent/internal/models"
	"testing"
)

type MockUserRepo struct {
	OnCreate     func(u *models.User) (int, error)
	OnGetByEmail func(email string) (*models.User, error)
}

func (m *MockUserRepo) Create(ctx context.Context, u *models.User) (int, error) {
	return m.OnCreate(u)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.OnGetByEmail(email)
}

func TestAuthService_Register(t *testing.T) {
	t.Run("успешная регистрация", func(t *testing.T) {
		mock := &MockUserRepo{
			OnCreate: func(u *models.User) (int, error) {
				return 1, nil
			},
		}
		svc := NewAuthService(mock)

		id, err := svc.Register(context.Background(), &models.User{Email: "test@test.com"})

		if err != nil {
			t.Errorf("ожидалось отсутствие ошибки, получено: %v", err)
		}
		if id != 1 {
			t.Errorf("ожидался ID 1, получено: %d", id)
		}
	})

	t.Run("ошибка базы данных", func(t *testing.T) {
		mock := &MockUserRepo{
			OnCreate: func(u *models.User) (int, error) {
				return 0, errors.New("db connection error")
			},
		}
		svc := NewAuthService(mock)

		_, err := svc.Register(context.Background(), &models.User{})

		if err == nil {
			t.Error("ожидалась ошибка, но получен nil")
		}
	})
}
