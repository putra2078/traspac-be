package taskCard

import (
	"hrm-app/internal/domain/labels"
	"time"
)

type TaskCard struct {
	ID        uint                   `json:"id" gorm:"primarykey"`
	TaskTabID uint                   `json:"task_tab_id"`
	Name      string                 `json:"name"`
	Content   string                 `json:"content"`
	Comment   string                 `json:"comment"`
	Date      string                 `json:"date"`
	Status    bool                   `json:"status"`
	Labels    []labels.TaskCardLabel `json:"labels" gorm:"foreignKey:TaskCardID"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
