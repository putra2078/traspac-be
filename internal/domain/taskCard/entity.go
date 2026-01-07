package taskCard

import (
	"hrm-app/internal/domain/labels"
	"hrm-app/internal/domain/taskCardComment"
	"hrm-app/internal/domain/taskCardUsers"
	"time"
)

type TaskCard struct {
	ID        uint                              `json:"id" gorm:"primarykey"`
	TaskTabID uint                              `json:"task_tab_id"`
	Name      string                            `json:"name"`
	Content   string                            `json:"content"`
	Date      string                            `json:"date"`
	Status    bool                              `json:"status"`
	Labels    []labels.TaskCardLabel            `json:"labels" gorm:"foreignKey:TaskCardID"`
	Comments  []taskCardComment.TaskCardComment `json:"comments" gorm:"foreignKey:TaskCardID"`
	Members   []taskCardUsers.TaskCardUsers     `json:"members" gorm:"foreignKey:TaskCardID"`
	CreatedAt time.Time                         `json:"created_at"`
	UpdatedAt time.Time                         `json:"updated_at"`
}
