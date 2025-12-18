package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestRequestIDMiddlewareAddsHeaderAndContext(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.Use(RequestIDMiddleware())
    r.GET("/test", func(c *gin.Context) {
        rid := c.GetString(RequestIDKey)
        if rid == "" {
            t.Fatalf("expected request id in context")
        }
        if got := c.Writer.Header().Get("X-Request-ID"); got == "" {
            t.Fatalf("expected X-Request-ID header to be set")
        }
        c.JSON(http.StatusOK, gin.H{"request_id": rid})
    })

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200 OK, got %d", w.Code)
    }
}
