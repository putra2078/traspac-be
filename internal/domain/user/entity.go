package user

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Photo       string `json:"photo"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=15"`
	Gender      string `json:"gender" binding:"required,oneof=male female other"`
	Address     string `json:"address"`
	BirthDate   string `json:"birth_date" binding:"required"` // Receive as string to parse manually or let binding handle if format is strict
}
