package handler

import (
	"goevent/internal/handler/auth"
	"goevent/internal/handler/event"
	"goevent/internal/service"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	handler *event.Handler
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{
		handler: event.NewHandler(s),
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	h.handler.CreateEvent(c)
}

func (h *EventHandler) GetAllEvents(c *gin.Context) {
	h.handler.ListEvents(c)
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	h.handler.GetEvent(c)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	h.handler.UpdateEvent(c)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	h.handler.DeleteEvent(c)
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c *gin.Context) {
	auth.RegisterUser(c)
}

func (h *AuthHandler) Login(c *gin.Context) {
	auth.Login(c)
}
