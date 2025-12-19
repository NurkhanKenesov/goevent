package invite

import (
	"goevent/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	authService *service.AuthService
}

func NewHandler(s *Service, as *service.AuthService) *Handler {
	return &Handler{
		service:     s,
		authService: as,
	}
}

func (h *Handler) InviteUser(c *gin.Context) {
	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req struct {
		InviteeID int64 `json:"invitee_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email"})
		return
	}

	u, err := h.authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	err = h.service.InviteUser(c.Request.Context(), int64(u.ID), req.InviteeID, eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation sent"})
}

func (h *Handler) RespondInvitation(c *gin.Context) {
	invitationIDStr := c.Param("id")
	invitationID, err := strconv.ParseInt(invitationIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invitation id"})
		return
	}

	var req struct {
		Status Status `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email"})
		return
	}

	u, err := h.authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	err = h.service.RespondInvitation(c.Request.Context(), int64(u.ID), invitationID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "response recorded"})
}

func (h *Handler) GetMyEvents(c *gin.Context) {
	emailVal, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid email"})
		return
	}

	u, err := h.authService.GetByEmail(c.Request.Context(), email)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	events, err := h.service.GetMyEvents(c.Request.Context(), int64(u.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
