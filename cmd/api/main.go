package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"goevent/internal/handler"
	"goevent/internal/middleware"
	"goevent/internal/repository"
	"goevent/internal/service"
)

func main() {
	// Подключение к SQLite (файловая база, не требует установки PostgreSQL)
	db, err := sqlx.Connect("sqlite3", "./goevent.db")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Создаём таблицы (простая версия)
	db.MustExec(`
	    CREATE TABLE IF NOT EXISTS users (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        username TEXT UNIQUE NOT NULL,
	        email TEXT UNIQUE NOT NULL,
	        password TEXT NOT NULL,
	        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	    )
	`)

	db.MustExec(`
	    CREATE TABLE IF NOT EXISTS events (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        title TEXT NOT NULL,
	        description TEXT,
	        date DATETIME NOT NULL,
	        location TEXT NOT NULL,
	        creator_id INTEGER NOT NULL,
	        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	        FOREIGN KEY (creator_id) REFERENCES users(id)
	    )
	`)

	db.MustExec(`
	    CREATE TABLE IF NOT EXISTS invitations (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        event_id INTEGER NOT NULL,
	        invitee_id INTEGER NOT NULL,
	        inviter_id INTEGER NOT NULL,
	        status TEXT NOT NULL DEFAULT 'pending',
	        message TEXT,
	        sent_at DATETIME NOT NULL,
	        responded_at DATETIME,
	        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	        FOREIGN KEY (event_id) REFERENCES events(id),
	        FOREIGN KEY (invitee_id) REFERENCES users(id),
	        FOREIGN KEY (inviter_id) REFERENCES users(id)
	    )
	`)

	log.Println("✅ Database connected and initialized")

	// Инициализация зависимостей
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)

	authService := service.NewAuthService(userRepo, "your-secret-key")
	eventService := service.NewEventService(eventRepo)
	invitationService := service.NewInvitationService(invitationRepo, eventRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	eventHandler := handler.NewEventHandler(eventService)
	invitationHandler := handler.NewInvitationHandler(invitationService)

	r := gin.Default()

	// Public routes - аутентификация
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.GET("/profile", authHandler.GetProfile)

		// Event routes
		events := api.Group("/events")
		{
			events.POST("", eventHandler.CreateEvent)
			events.GET("", eventHandler.GetMyEvents)
			events.GET("/:id", eventHandler.GetEvent)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
		}

		// Invitation routes
		invitations := api.Group("/invitations")
		{
			invitations.POST("", invitationHandler.CreateInvitation)
			invitations.GET("", invitationHandler.GetMyInvitations)
			invitations.GET("/sent", invitationHandler.GetSentInvitations)
			invitations.GET("/:id", invitationHandler.GetInvitation)
			invitations.GET("/:id/details", invitationHandler.GetInvitationDetails)
			invitations.PUT("/:id/respond", invitationHandler.RespondToInvitation)
			invitations.DELETE("/:id", invitationHandler.CancelInvitation)
		}

		// Event-specific invitation routes
		events.Group("/:eventId").GET("/invitations", invitationHandler.GetEventInvitations)
	}

	// Public event routes (read-only)
	r.GET("/events", eventHandler.GetAllEvents)
	r.GET("/events/:id", eventHandler.GetEvent)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})

	log.Println("🚀 Server starting on :4000")
	r.Run(":4000")
}
