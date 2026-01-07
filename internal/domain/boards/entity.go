package boards

import (
	"hrm-app/internal/domain/taskTab"
	"time"
)

type Boards struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	WorkspaceID uint              `json:"workspace_id"`
	CreatedBy   uint              `json:"created_by"`
	Name        string            `json:"name"`
	Images      string            `json:"images"`
	TaskTabs    []taskTab.TaskTab `json:"task_tabs" gorm:"foreignKey:BoardID"`
	TaskCards   interface{}       `json:"task_cards" gorm:"-"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
