package taskTab

import (
	"time"
)

type TaskTab struct {
	ID        		uint      `json:"id" gorm:"primaryKey"`
	BoardID 		uint      `json:"board_id"`
	Position 		int 	  `json:"position"`
	Name      		string    `json:"name"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`
}