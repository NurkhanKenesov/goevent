package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	h := &AuthHandler{
		jwtSecret: "test_secret",
	}

	r.POST("/login", h.Login)

	t.Run("invalid json", func(t *testing.T) {
		body := []byte(`{"email": "wrong-json"`)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
func TestAuthHandler_Register_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		inputBody      interface{} // Изменили на interface{}, чтобы передавать строки
		expectedStatus int
	}{
		{
			name: "Success Registration",
			inputBody: map[string]interface{}{
				"username": "newuser",
				"email":    "new@test.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			inputBody:      "{broken-json}", // Теперь отправляем реально сломанный JSON
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := &AuthHandler{authService: &MockAuthService{}}
			r.POST("/register", h.RegisterUser)

			var buf bytes.Buffer
			// Логика кодирования: если это строка, пишем как есть, если мапа — кодируем в JSON
			if s, ok := tt.inputBody.(string); ok {
				buf.WriteString(s)
			} else {
				json.NewEncoder(&buf).Encode(tt.inputBody)
			}

			req := httptest.NewRequest("POST", "/register", &buf)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected %d, got %d. Body: %s", tt.name, tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}
func TestAuthHandler_Profile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Profile", func(t *testing.T) {
		r := gin.New()
		h := &AuthHandler{authService: &MockAuthService{}}

		r.GET("/profile", func(c *gin.Context) {
			c.Set("email", "test@test.com")
			h.Profile(c)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/profile", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("No email in context", func(t *testing.T) {
		r := gin.New()
		h := &AuthHandler{authService: &MockAuthService{}}

		r.GET("/profile", h.Profile)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/profile", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
}
