package service

import (
	"context"
	"goevent/internal/models"
	"goevent/internal/repository"
)

type AuthService struct {
	repo *repository.UserRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.NewUserRepo(),
	}
}

func (s *AuthService) Register(ctx context.Context, u *models.User) (int, error) {
	return s.repo.Create(ctx, u)
}

func (s *AuthService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
