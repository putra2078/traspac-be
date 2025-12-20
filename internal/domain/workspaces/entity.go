package workspaces

import (
	"time"
)

type Workspace struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedBy uint      `json:"created_by"`
	PassCode  string    `json:"pass_code"`
	Name      string    `json:"name"`
	Privacy   string    `json:"privacy"`
	JoinLink  string    `json:"join_link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
