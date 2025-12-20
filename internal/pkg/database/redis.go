package database

import (
	"context"
	"fmt"
	"log"

	"hrm-app/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRedis(cfg *config.Config) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
