package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hrm-app/config"
	"hrm-app/internal/app"
	"hrm-app/internal/middleware"
	"hrm-app/internal/pkg/database"
	rmqConfig "hrm-app/internal/pkg/rabbitmq/config"
	rmqConnection "hrm-app/internal/pkg/rabbitmq/connection"
	rmqManager "hrm-app/internal/pkg/rabbitmq/manager"

	"github.com/gin-gonic/gin"
)

// func main() {
// 	// Load configuration
// 	cfg := config.LoadConfig()

// 	// Set Gin mode
// 	mode := os.Getenv("GIN_MODE")
// 	if mode == "" {
// 		mode = "debug"
// 	}
// 	gin.SetMode(mode)

// 	// --- Connect to PostgreSQL ---
// 	database.ConnectDatabase(cfg)
// 	log.Println("[INFO] PostgreSQL connected successfully")

// 	// --- Connect to Redis ---
// 	database.ConnectRedis(cfg)
// 	log.Println("[INFO] Redis connected successfully")

// 	// --- Setup Gin Router ---
// 	r, hub := app.SetupRouter(cfg)
// 	r.Use(gin.Recovery())

// 	// --- Server Configuration ---
// 	port := fmt.Sprintf(":%d", cfg.Server.Port)
// 	srv := &http.Server{
// 		Addr:              port,
// 		Handler:           r,
// 		ReadHeaderTimeout: 5 * time.Second,
// 		ReadTimeout:       15 * time.Second,
// 		WriteTimeout:      15 * time.Second,
// 		IdleTimeout:       60 * time.Second,
// 	}

// 	// Jalankan server di goroutine
// 	go func() {
// 		log.Printf("[INFO] Server is running on %s\n", port)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("[ERROR] Server failed: %v", err)
// 		}
// 	}()

// 	// --- Graceful Shutdown ---
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	log.Println("[INFO] Shutting down server...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
// 	}

// 	// Shutdown Hub and Workers
// 	if hub != nil {
// 		hub.Shutdown()
// 	}

// 	log.Println("[INFO] Server exited gracefully")
// }

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug"
	}
	gin.SetMode(mode)

	// --- Connect to PostgreSQL ---
	database.ConnectDatabase(cfg)
	log.Println("[INFO] PostgreSQL connected successfully")

	// --- Connect to Redis ---
	database.ConnectRedis(cfg)
	log.Println("[INFO] Redis connected successfully")

	// --- Connect to RabbitMQ ---
	rmqConn, err := rmqConnection.New(rmqConfig.RabbitURL)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to RabbitMQ: %v", err)
	}
	defer rmqConn.Close()
	log.Println("[INFO] RabbitMQ connected successfully")

	// --- Initialize RabbitMQ Manager ---
	channelManager := rmqManager.NewChannelManager(rmqConn)
	rateLimiter := middleware.NewRateLimiter()
	log.Println("[INFO] RabbitMQ Channel Manager initialized")

	// --- Setup Gin Router ---
	r, hub := app.SetupRouter(cfg, channelManager, rateLimiter)
	r.Use(gin.Recovery())

	// --- Server Configuration ---
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Jalankan server di goroutine
	go func() {
		log.Printf("[INFO] Server is running on %s\n", port)
		log.Printf("[INFO] RabbitMQ Management UI: http://localhost:15672")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}

	// Shutdown Hub and Workers
	if hub != nil {
		hub.Shutdown()
	}

	// Close all RabbitMQ channels
	log.Println("[INFO] Closing RabbitMQ channels...")
	// channelManager will auto-cleanup when connection closes

	log.Println("[INFO] Server exited gracefully")
}
