package room_messages

import (
	"hrm-app/internal/domain/user"
	"time"
)

type RoomMessage struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         *uint      `json:"user_id"`
	User           *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RoomID         uint       `json:"room_id"`
	MessageText    string     `json:"message_text"`
	MessageContent string     `json:"message_content"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (RoomMessage) TableName() string {
	return "room_messages"
}
