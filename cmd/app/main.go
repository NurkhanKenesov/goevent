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
	"goevent/internal/repository"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	jwtSecret := "mysecretkey"

	eventRepo := event.NewRepository(database)
	userRepo := repository.NewUserRepo()

	eventService := event.NewService(eventRepo, userRepo)
	eventHandler := eventhandler.NewHandler(eventService)

	inviteRepo := invite.NewRepository(database)
	inviteService := invite.NewService(inviteRepo, eventRepo)
	inviteHandler := invite.NewHandler(inviteService)

	r := gin.Default()

	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.Login)

	authMiddleware := auth.AuthMiddleware(jwtSecret)

	events := r.Group("/events")
	{
		events.GET("", eventHandler.ListEvents)
		events.GET("/:id", eventHandler.GetEvent)

		events.Use(authMiddleware)
		{
			events.POST("", eventHandler.CreateEvent)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
		}
	}

	r.POST("/events/:id/invite", authMiddleware, inviteHandler.InviteUser)
	r.POST("/invitations/:id/respond", authMiddleware, inviteHandler.RespondInvitation)
	r.GET("/my-events", authMiddleware, inviteHandler.GetMyEvents)

	r.Run(":3000")
}
