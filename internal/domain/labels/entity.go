package labels

import (
	"time"
)

type TaskCardLabel struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	TaskCardID  uint        `json:"task_card_id"`
	Title       string      `json:"title"`
	Color       string      `json:"color"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}