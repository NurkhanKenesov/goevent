package service

import (
	"context"
	"goevent/internal/models"
)

type AuthRepository interface {
	Create(ctx context.Context, u *models.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, u *models.User) (int64, error) {
	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (s *AuthService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
