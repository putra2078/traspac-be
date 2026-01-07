package boardsUsers

import (
	"hrm-app/internal/domain/user"
	"time"
)

type BoardsUsers struct {
	ID        uint       `json:"id"`
	BoardID   uint       `json:"board_id"`
	UserID    uint       `json:"user_id"`
	User      *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
