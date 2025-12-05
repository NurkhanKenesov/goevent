package main

import (
	"fmt"
	"goevent/internal/config"
	"goevent/internal/db"
	"goevent/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	fmt.Println("DB connected")

	r := gin.Default()

	r.GET("/ping", handler.Ping)
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.Login)

	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware(cfg.JWTSecret))
	auth.GET("/profile", handler.Profile)

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
