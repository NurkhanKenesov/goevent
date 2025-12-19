package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"goevent/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockAuthService struct{}

func (m *MockAuthService) Register(ctx context.Context, u *models.User) (int64, error) {
	return 1, nil
}

func (m *MockAuthService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, nil
	}
	return &models.User{ID: 1, Email: email, Password: "hashed_password"}, nil
}

func TestAuthHandler_Login_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputBody      map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "Invalid JSON",
			inputBody:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty email",
			inputBody: map[string]interface{}{
				"email":    "",
				"password": "123",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			h := &AuthHandler{
				authService: &MockAuthService{},
				jwtSecret:   "secret",
			}
			r.POST("/login", h.Login)

			var buf bytes.Buffer
			if tt.inputBody != nil {
				json.NewEncoder(&buf).Encode(tt.inputBody)
			}

			req := httptest.NewRequest("POST", "/login", &buf)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}
func TestAuthHandler_Login_NegativeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputBody      map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "User not found (401)",
			inputBody: map[string]interface{}{
				"email":    "nonexistent@test.com",
				"password": "any_password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid credentials",
		},
		{
			name: "Wrong password (401)",
			inputBody: map[string]interface{}{
				"email":    "test@test.com",
				"password": "wrong_password",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := &AuthHandler{
				authService: &MockAuthService{},
				jwtSecret:   "secret",
			}
			r.POST("/login", h.Login)

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(tt.inputBody)

			req := httptest.NewRequest("POST", "/login", &buf)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			if response["error"] != tt.expectedError {
				t.Errorf("%s: expected error message '%s', got '%s'", tt.name, tt.expectedError, response["error"])
			}
		})
	}
}
