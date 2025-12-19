package invite

import (
	"goevent/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInviteHandler_All(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := NewService(&mockInviteRepo{}, &mockEventRepo{})
	h := NewHandler(mockSvc, &service.AuthService{})

	t.Run("Invalid Event ID (400)", func(t *testing.T) {
		r := gin.New()
		r.POST("/events/:id/invite", h.InviteUser)

		req := httptest.NewRequest("POST", "/events/abc/invite", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Missing Email in Context (401)", func(t *testing.T) {
		r := gin.New()
		r.POST("/events/:id/invite", h.InviteUser)

		body := `{"invitee_id": 10}`
		req := httptest.NewRequest("POST", "/events/1/invite", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
