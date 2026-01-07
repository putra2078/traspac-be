package taskCardUsers

import (
	"hrm-app/internal/domain/user"
	"time"
)

type TaskCardUsers struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	TaskCardID uint      `json:"task_card_id"`
	UserID     uint      `json:"user_id"`
	User       user.User `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
