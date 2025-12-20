package servers

import "time"

type Server struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedBy uint      `json:"created_by"`
	Name      string    `json:"name"`
	Privacy   string    `json:"privacy"`
	PassCode  string    `json:"pass_code"`
	LinkJoin  string    `json:"link_join"`
	CreatedAt time.Time `json:"created_at"`
	Updatedat time.Time `json:"updated_at"`
}
