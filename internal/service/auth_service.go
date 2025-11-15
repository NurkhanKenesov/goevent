package service

import (
	"errors"
	"time"

	"goevent/internal/models"
	"goevent/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Проверяем, существует ли email
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем в БД
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Генерируем JWT токен
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Генерируем JWT токен
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := int(claims["user_id"].(float64))
	return s.userRepo.GetUserByID(userID)
}
