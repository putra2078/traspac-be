package utils

import (
	"context"
	"fmt"
	"time"

	"hrm-app/internal/pkg/database"
)

var ctx = context.Background()

func SetSession(userID uint, token string, exp time.Duration) error {
	key := fmt.Sprintf("session:%d:%s", userID, token)
	return database.RDB.Set(ctx, key, "active", exp).Err()
}

func GetSession(userID uint, token string) (string, error) {
	key := fmt.Sprintf("session:%d:%s", userID, token)
	return database.RDB.Get(ctx, key).Result()
}

func DeleteSession(userID uint, token string) error {
	key := fmt.Sprintf("session:%d:%s", userID, token)
	return database.RDB.Del(ctx, key).Err()
}

func ExtendSession(userID uint, token string, exp time.Duration) error {
	key := fmt.Sprintf("session:%d:%s", userID, token)
	return database.RDB.Expire(ctx, key, exp).Err()
}
