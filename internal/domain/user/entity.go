package user

import (
	"mime/multipart"
	"time"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"-"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type RegisterRequest struct {
	Username    string                `json:"username" form:"username" binding:"required,min=3,max=50"`
	Email       string                `json:"email" form:"email" binding:"required,email"`
	Password    string                `json:"password" form:"password" binding:"required,min=8"`
	Name        string                `json:"name" form:"name" binding:"required,min=3,max=100"`
	Photo       string                `json:"photo"`               // Kept for backward compatibility or direct URL
	PhotoFile   *multipart.FileHeader `json:"-" form:"photo_file"` // New field for file upload
}

type UpdateRequest struct {
	Name        string                `json:"name" form:"name"` // Optional updates
	Username    string                `json:"username" form:"username"`
	Photo       string                `json:"photo"`
	PhotoFile   *multipart.FileHeader `json:"-" form:"photo_file"`
}
