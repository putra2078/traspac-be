package room_chats

import "time"

type RoomsChats struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	WorkspaceID uint      `json:"workspace_id"`
	Name        string    `json:"name"`
	CreatedBy   uint      `json:"created_by"`
	PassCode    string    `json:"passcode"`
	LinkJoin    string    `json:"link_join"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
