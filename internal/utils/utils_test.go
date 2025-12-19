package utils

import (
	"testing"
	"time"
)

func TestPasswordHashing(t *testing.T) {
	password := "my_super_secret_123"

	// Тестируем хеширование
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Тестируем сравнение
	if !CompareHashAndPassword(hashed, password) {
		t.Error("Password comparison failed for correct password")
	}

	if CompareHashAndPassword(hashed, "wrong_password") {
		t.Error("Password comparison succeeded for wrong password")
	}
}

func TestJWTToken(t *testing.T) {
	secret := "test_secret_key"
	email := "user@example.com"
	ttl := 1 * time.Hour

	// 1. Тестируем создание
	token, err := GenerateToken(email, secret, ttl)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 2. Тестируем парсинг (валидацию)
	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to parse valid token: %v", err)
	}

	if claims["email"] != email {
		t.Errorf("Expected email %s, got %v", email, claims["email"])
	}

	// 3. Тестируем просроченный токен (пункт чек-листа: Token expiration)
	expiredToken, _ := GenerateToken(email, secret, -1*time.Minute)
	_, err = ParseToken(expiredToken, secret)
	if err == nil {
		t.Error("Expected error for expired token, but got nil")
	}
}
