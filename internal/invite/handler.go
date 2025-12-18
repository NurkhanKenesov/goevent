package invite

import (
	"encoding/json"
	"goevent/internal/errors"
	"goevent/internal/idempotency"
	"goevent/internal/logging"
	"goevent/internal/middleware"
	"goevent/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var authService = service.NewAuthService()

type Handler struct {
	service *Service
	store   *idempotency.Store
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s, store: idempotency.Default}
}

// normalizeIdempotencyKey validates and normalizes the Idempotency-Key header.
func normalizeIdempotencyKey(c *gin.Context) (string, *errors.AppError) {
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if key == "" {
		return "", errors.ErrValidation("Idempotency-Key header is required")
	}
	if len(key) > 255 {
		return "", errors.ErrValidation("Idempotency-Key must be <= 255 characters")
	}
	return key, nil
}

func (h *Handler) InviteUser(c *gin.Context) {
	rid := middleware.FromContext(c.Request.Context())

	// Extract and validate Idempotency-Key
	idempotencyKey, appErr := normalizeIdempotencyKey(c)
	if appErr != nil {
		logging.Error("invite.validation.failed", map[string]interface{}{
			"request_id": rid,
			"error":      appErr.Message,
		})
		statusCode, errResp := errors.ToHTTP(appErr)
		c.JSON(statusCode, errResp)
		return
	}

	// Check cache first
	if cachedStatus, cachedBody, found := h.store.Get(idempotencyKey); found {
		logging.Info("invite.cache.hit", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
			"cached_status":   cachedStatus,
		})
		c.JSON(cachedStatus, gin.H{"cached": true, "data": json.RawMessage(cachedBody)})
		return
	}

	// Try to acquire lock to prevent concurrent execution
	if !h.store.TryLock(idempotencyKey) {
		logging.Warn("invite.concurrent_request", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
		})
		// Another request with the same key is in flight
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate_request", "message": "request already in progress"})
		return
	}
	defer h.store.Unlock(idempotencyKey)

	// Check cache again after lock (double-check pattern)
	if cachedStatus, cachedBody, found := h.store.Get(idempotencyKey); found {
		logging.Info("invite.cache.hit.after_lock", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
			"cached_status":   cachedStatus,
		})
		c.JSON(cachedStatus, gin.H{"cached": true, "data": json.RawMessage(cachedBody)})
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrValidation("invalid event id"))
		c.JSON(statusCode, errResp)
		return
	}

	var req struct {
		InviteeID int64 `json:"invitee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrValidation(err.Error()))
		c.JSON(statusCode, errResp)
		return
	}

	emailVal, exists := c.Get("email")
	if !exists {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("unauthorized"))
		c.JSON(statusCode, errResp)
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal("invalid email"))
		c.JSON(statusCode, errResp)
		return
	}

	u, err := authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("user not found"))
		c.JSON(statusCode, errResp)
		return
	}

	err = h.service.InviteUser(c.Request.Context(), int64(u.ID), req.InviteeID, eventID)
	if err != nil {
		logging.Error("invite.service.error", map[string]interface{}{
			"request_id": rid,
			"error":      err.Error(),
		})
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal(err.Error()))
		c.JSON(statusCode, errResp)
		return
	}

	// Marshal and cache successful response
	respBody := gin.H{"message": "invitation sent"}
	respBytes, _ := json.Marshal(respBody)
	h.store.Set(idempotencyKey, http.StatusOK, respBytes)

	logging.Info("invite.executed", map[string]interface{}{
		"request_id":      rid,
		"idempotency_key": idempotencyKey,
		"event_id":        eventID,
		"invitee_id":      req.InviteeID,
	})

	c.JSON(http.StatusOK, respBody)
}

func (h *Handler) RespondInvitation(c *gin.Context) {
	rid := middleware.FromContext(c.Request.Context())

	// Extract and validate Idempotency-Key
	idempotencyKey, appErr := normalizeIdempotencyKey(c)
	if appErr != nil {
		logging.Error("respond.validation.failed", map[string]interface{}{
			"request_id": rid,
			"error":      appErr.Message,
		})
		statusCode, errResp := errors.ToHTTP(appErr)
		c.JSON(statusCode, errResp)
		return
	}

	// Check cache first
	if cachedStatus, cachedBody, found := h.store.Get(idempotencyKey); found {
		logging.Info("respond.cache.hit", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
			"cached_status":   cachedStatus,
		})
		c.JSON(cachedStatus, gin.H{"cached": true, "data": json.RawMessage(cachedBody)})
		return
	}

	// Try to acquire lock to prevent concurrent execution
	if !h.store.TryLock(idempotencyKey) {
		logging.Warn("respond.concurrent_request", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
		})
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate_request", "message": "request already in progress"})
		return
	}
	defer h.store.Unlock(idempotencyKey)

	// Check cache again after lock (double-check pattern)
	if cachedStatus, cachedBody, found := h.store.Get(idempotencyKey); found {
		logging.Info("respond.cache.hit.after_lock", map[string]interface{}{
			"request_id":      rid,
			"idempotency_key": idempotencyKey,
			"cached_status":   cachedStatus,
		})
		c.JSON(cachedStatus, gin.H{"cached": true, "data": json.RawMessage(cachedBody)})
		return
	}

	invitationIDStr := c.Param("id")
	invitationID, err := strconv.ParseInt(invitationIDStr, 10, 64)
	if err != nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrValidation("invalid invitation id"))
		c.JSON(statusCode, errResp)
		return
	}

	var req struct {
		Status Status `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrValidation(err.Error()))
		c.JSON(statusCode, errResp)
		return
	}

	emailVal, exists := c.Get("email")
	if !exists {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("unauthorized"))
		c.JSON(statusCode, errResp)
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal("invalid email"))
		c.JSON(statusCode, errResp)
		return
	}

	u, err := authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("user not found"))
		c.JSON(statusCode, errResp)
		return
	}

	err = h.service.RespondInvitation(c.Request.Context(), int64(u.ID), invitationID, req.Status)
	if err != nil {
		logging.Error("respond.service.error", map[string]interface{}{
			"request_id": rid,
			"error":      err.Error(),
		})
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal(err.Error()))
		c.JSON(statusCode, errResp)
		return
	}

	// Marshal and cache successful response
	respBody := gin.H{"message": "response recorded"}
	respBytes, _ := json.Marshal(respBody)
	h.store.Set(idempotencyKey, http.StatusOK, respBytes)

	logging.Info("respond.executed", map[string]interface{}{
		"request_id":      rid,
		"idempotency_key": idempotencyKey,
		"invitation_id":   invitationID,
		"status":          req.Status,
	})

	c.JSON(http.StatusOK, respBody)
}

func (h *Handler) GetMyEvents(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("unauthorized"))
		c.JSON(statusCode, errResp)
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal("invalid email"))
		c.JSON(statusCode, errResp)
		return
	}

	u, err := authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrUnauthorized("user not found"))
		c.JSON(statusCode, errResp)
		return
	}

	events, err := h.service.GetMyEvents(c.Request.Context(), int64(u.ID))
	if err != nil {
		statusCode, errResp := errors.ToHTTP(errors.ErrInternal(err.Error()))
		c.JSON(statusCode, errResp)
		return
	}

	c.JSON(http.StatusOK, events)
}
