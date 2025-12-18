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

func (s *AuthService) Register(
	ctx context.Context,
	u *models.User,
) (int64, error) {

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (s *AuthService) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
