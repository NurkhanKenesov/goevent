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

	// Создаём таблицу users (простая версия)
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

	log.Println("✅ Database connected and initialized")

	// Инициализация зависимостей
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "your-secret-key")
	authHandler := handler.NewAuthHandler(authService)

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
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})

	log.Println("🚀 Server starting on :4000")
	r.Run(":4000")
}
