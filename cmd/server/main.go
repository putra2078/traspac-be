// package main

// import (
// 	"fmt"
// 	"hrm-app/config"
// 	"hrm-app/internal/app"
// 	"hrm-app/internal/domain/contact"
// 	"hrm-app/internal/domain/employee"
// 	"hrm-app/internal/domain/manager"
// 	"hrm-app/internal/domain/user"
// 	"hrm-app/internal/pkg/database"
// )

// func main() {
// 	cfg := config.LoadConfig()

// 	database.ConnectDatabase(cfg)
// 	database.ConnectRedis(cfg)

// 	r := app.SetupRouter(cfg)
// 	port := fmt.Sprintf(":%d", cfg.Server.Port)
// 	// Auto migrate database schemas
// 	// ensure contact table exists as we now create contacts in a transaction
// 	database.DB.AutoMigrate(&employee.Employee{}, &contact.Contact{}, &user.User{}, &manager.Manager{})
// 	r.Run(port)
// }

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

	// "hrm-app/internal/domain/contact"
	// "hrm-app/internal/domain/department"
	// "hrm-app/internal/domain/employee"
	// "hrm-app/internal/domain/manager"
	// "hrm-app/internal/domain/user"
	// "hrm-app/internal/middleware"
	"hrm-app/internal/pkg/database"

	"github.com/gin-gonic/gin"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode (default: release untuk production)
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug" // default ke production-safe mode
	}
	gin.SetMode(mode)

	// --- Initialize Prometheus Metrics ---
	// middleware.InitPrometheus()
	// log.Println("[INFO] Prometheus metrics initialized")

	// --- Connect to PostgreSQL ---
	database.ConnectDatabase(cfg) // function ini sudah handle error & logging internal
	log.Println("[INFO] PostgreSQL connected successfully")

	// --- Connect to Redis ---
	// Optional: hanya jalankan jika kamu punya file redis.go
	database.ConnectRedis(cfg)
	log.Println("[INFO] Redis connected successfully")

	// --- Auto Migration (hanya di mode debug) ---
	// if gin.Mode() == gin.DebugMode {
	// 	log.Println("[INFO] Running auto migration (debug mode only)...")
	// 	err := database.DB.AutoMigrate(
	// 		&employee.Employee{},
	// 		&contact.Contact{},
	// 		&user.User{},
	// 		&manager.Manager{},
	// 		&department.Department{},
	// 	)
	// 	if err != nil {
	// 		log.Fatalf("[ERROR] Migration failed: %v", err)
	// 	}
	// 	log.Println("[INFO] Auto migration completed successfully")
	// }

	// --- Setup Gin Router ---
	r := app.SetupRouter(cfg)
	r.Use(gin.Recovery()) // recover dari panic
	// --- Setup Prometheus Metrics Endpoint ---
	// r.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

	// --- Server Configuration ---
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,  // ⏳ Timeout untuk membaca header
		ReadTimeout:       15 * time.Second, // batas waktu baca body
		WriteTimeout:      15 * time.Second, // batas waktu tulis respons
		IdleTimeout:       60 * time.Second, // waktu idle maksimum
	}

	// Jalankan server di goroutine agar bisa shutdown dengan elegan
	go func() {
		log.Printf("[INFO] Server is running on %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Server failed: %v", err)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited gracefully")
}
