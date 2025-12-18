package eventhandler

import (
	"errors"
	"net/http"
	"strconv"

	"goevent/internal/event"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *event.Service
}

func NewHandler(s *event.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var e event.Event
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	emailRaw, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email := emailRaw.(string)

	if err := h.service.CreateEvent(c.Request.Context(), &e, email); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, e)
}

func (h *Handler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	e, err := h.service.GetEvent(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var e event.Event
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	emailRaw, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email := emailRaw.(string)

	if err := h.service.UpdateEvent(
		c.Request.Context(),
		id,
		&e,
		email,
	); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, e)
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	emailRaw, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email := emailRaw.(string)

	if err := h.service.DeleteEvent(
		c.Request.Context(),
		id,
		email,
	); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListEvents(c *gin.Context) {
	events, err := h.service.ListEvents(c.Request.Context(), 10, 0)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, event.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, event.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
