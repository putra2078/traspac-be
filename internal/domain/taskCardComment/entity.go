package taskCardComment

import (
	"time"
)

type TaskCardComment struct {
	ID        int       `json:"id"`
	TaskCardID  int       `json:"task_card_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
