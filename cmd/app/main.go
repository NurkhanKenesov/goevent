package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"goevent/internal/config"
	"goevent/internal/db"
	"goevent/internal/event"
	"goevent/internal/handler/auth"
	eventhandler "goevent/internal/handler/event"
	"goevent/internal/invite"
	"goevent/internal/middleware"
	"goevent/internal/logging"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}
	defer database.Close()

	jwtSecret := cfg.JWTSecret

	eventRepo := event.NewRepository(database)
	eventService := event.NewService(eventRepo)
	eventHandler := eventhandler.NewHandler(eventService)

	inviteRepo := invite.NewRepository(database)
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()
	
	inviteService := invite.NewService(workerCtx, inviteRepo, eventRepo)
	inviteHandler := invite.NewHandler(inviteService)

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware())

	r.POST("/register", auth.RegisterUser)
	r.POST("/login", auth.Login)

	r.POST("/events", eventHandler.CreateEvent)
	r.GET("/events", eventHandler.ListEvents)
	r.GET("/events/:id", eventHandler.GetEvent)
	r.PUT("/events/:id", eventHandler.UpdateEvent)
	r.DELETE("/events/:id", eventHandler.DeleteEvent)

	authMiddleware := auth.AuthMiddleware(jwtSecret)
	r.POST("/events/:id/invite", authMiddleware, inviteHandler.InviteUser)
	r.POST("/invitations/:id/respond", authMiddleware, inviteHandler.RespondInvitation)
	r.GET("/my-events", authMiddleware, inviteHandler.GetMyEvents)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logging.Info("server.start", map[string]interface{}{"addr": server.Addr})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		logging.Info("server.signal_received", map[string]interface{}{"signal": "SIGINT/SIGTERM"})
	case err := <-serverErrors:
		log.Fatal("Server error: ", err)
	}

	logging.Info("server.shutting_down", nil)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cancelWorkers()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logging.Error("server.shutdown_error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	logging.Info("server.stopped", nil)
}
