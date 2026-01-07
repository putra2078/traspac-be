package roomUsers

import (
	"hrm-app/internal/domain/user"
	"time"
)

type RoomUsers struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	RoomID    uint       `json:"room_id"`
	UserID    uint       `json:"user_id"`
	User      *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
