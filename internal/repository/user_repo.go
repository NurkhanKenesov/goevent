package repository

import (
	"context"
	"goevent/internal/db"
	"goevent/internal/models"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { return &UserRepo{} }

func (r *UserRepo) Create(ctx context.Context, u *models.User) (int, error) {
	var id int
	err := db.Pool.QueryRow(ctx,
		"INSERT INTO users (username, email, password) VALUES ($1,$2,$3) RETURNING id",
		u.Username, u.Email, u.Password).Scan(&id)
	return id, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := db.Pool.QueryRow(ctx,
		"SELECT id, username, email, password, created_at FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	err := db.Pool.QueryRow(ctx,
		"SELECT id, username, email, created_at FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
