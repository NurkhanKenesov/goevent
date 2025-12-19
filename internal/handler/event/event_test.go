package event

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandler_GetEvent_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{} // Мы не инициализируем сервисы, так как тест упадет раньше их вызова

	r := gin.New()
	r.GET("/events/:id", h.GetEvent)

	// Тестируем передачу строки вместо ID
	req := httptest.NewRequest("GET", "/events/not-a-number", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
