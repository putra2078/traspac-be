package contact

import "time"

type Contact struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Photo       string    `json:"photo"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	BirthDate   time.Time `json:"birth_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
