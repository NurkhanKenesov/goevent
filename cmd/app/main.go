package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"goevent/internal/config"
	"goevent/internal/db"
	"goevent/internal/event"
	"goevent/internal/handler/auth"
	eventhandler "goevent/internal/handler/event"
	"goevent/internal/invite"
)

func main() {
	// Загружаем ENV
	cfg := config.Load()

	// Подключаемся к базе
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	jwtSecret := "mysecretkey"

	// Инициализация сервисов и хендлеров
	eventRepo := event.NewRepository(database)
	eventService := event.NewService(eventRepo)
	eventHandler := eventhandler.NewHandler(eventService)

	inviteRepo := invite.NewRepository(database)
	inviteService := invite.NewService(inviteRepo, eventRepo)
	inviteHandler := invite.NewHandler(inviteService)

	// Router
	r := gin.Default()

	// Auth
	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.Login)

	// Events
	r.POST("/events", eventHandler.CreateEvent)
	r.GET("/events", eventHandler.ListEvents)
	r.GET("/events/:id", eventHandler.GetEvent)
	r.PUT("/events/:id", eventHandler.UpdateEvent)
	r.DELETE("/events/:id", eventHandler.DeleteEvent)

	// Invitations
	authMiddleware := auth.AuthMiddleware(jwtSecret)
	r.POST("/events/:id/invite", authMiddleware, inviteHandler.InviteUser)
	r.POST("/invitations/:id/respond", authMiddleware, inviteHandler.RespondInvitation)
	r.GET("/my-events", authMiddleware, inviteHandler.GetMyEvents)

	// Запуск сервера
	r.Run(":3000")
}
