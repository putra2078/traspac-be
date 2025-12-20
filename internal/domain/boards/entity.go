package boards

import (
	"time"
)

type Boards struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	WorkspaceID uint        `json:"workspace_id"`
	Name        string      `json:"name"`
	Images      string      `json:"images"`
	TaskTabs    interface{} `json:"task_tabs" gorm:"-"`
	TaskCards   interface{} `json:"task_cards" gorm:"-"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
