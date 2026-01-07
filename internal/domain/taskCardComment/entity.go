package taskCardComment

import (
	"hrm-app/internal/domain/user"
	"time"
)

type TaskCardComment struct {
	ID         int       `json:"id"`
	TaskCardID int       `json:"task_card_id"`
	UserID     uint      `json:"user_id"`
	User       user.User `json:"user" gorm:"foreignKey:UserID"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
