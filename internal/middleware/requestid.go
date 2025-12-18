package middleware

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "time"

    "goevent/internal/logging"

    "github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func genRequestID() string {
    b := make([]byte, 12)
    if _, err := rand.Read(b); err != nil {
        // fallback to timestamp
        return "rid-" + time.Now().UTC().Format("20060102T150405Z07:00")
    }
    return hex.EncodeToString(b)
}

// RequestIDMiddleware ensures every request has a request ID, places it into
// the Gin context and the request's context.Context so it propagates to services.
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        rid := c.GetHeader("X-Request-ID")
        if rid == "" {
            rid = genRequestID()
        }
        // set header back for visibility
        c.Writer.Header().Set("X-Request-ID", rid)
        c.Set(RequestIDKey, rid)
        // also put into request.Context
        ctx := context.WithValue(c.Request.Context(), RequestIDKey, rid)
        c.Request = c.Request.WithContext(ctx)
        logging.Info("request.start", map[string]interface{}{"request_id": rid, "method": c.Request.Method, "path": c.FullPath()})
        start := time.Now()
        c.Next()
        duration := time.Since(start).Milliseconds()
        status := c.Writer.Status()
        logging.Info("request.end", map[string]interface{}{"request_id": rid, "method": c.Request.Method, "path": c.FullPath(), "status": status, "duration_ms": duration})
    }
}

// FromContext returns the request id string from ctx or empty.
func FromContext(ctx context.Context) string {
    if v := ctx.Value(RequestIDKey); v != nil {
        if s, ok := v.(string); ok {
            return s
        }
    }
    return ""
}
