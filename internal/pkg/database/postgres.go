package database

import (
	"fmt"
	"log"
	"time"

	"hrm-app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.Sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ✅ Setup connection pooling (ini yang bikin POST/GET lebih cepat & stabil)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB from gorm: %v", err)
	}

	// --- Pooling Configuration ---
	sqlDB.SetMaxOpenConns(20)           // Maksimum 20 koneksi aktif
	sqlDB.SetMaxIdleConns(10)           // Maksimum 10 koneksi idle (nganggur tapi siap dipakai)
	sqlDB.SetConnMaxLifetime(time.Hour) // Durasi maksimum 1 koneksi aktif = 1 jam

	// --- Optional Logging ---
	log.Println("✅ PostgreSQL database connected successfully")
	log.Printf("🌊 Connection pool configured (MaxOpenConns=%d, MaxIdleConns=%d)", 20, 10)

	DB = db
}
