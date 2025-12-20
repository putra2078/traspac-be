package user

import (
	"errors"

	"gorm.io/gorm"
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(user *User) error
	FindAll() ([]User, error)
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

type usersRepository struct{}

func NewRepository() Repository {
	return &usersRepository{}
}

func (r *usersRepository) Create(user *User) error {
	return database.DB.Create(user).Error
}

func (r *usersRepository) FindAll() ([]User, error) {
	var users []User
	err := database.DB.Find(&users).Error
	return users, err
}

func (r *usersRepository) FindByID(id uint) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error

	return &user, err
}

func (r *usersRepository) FindByEmail(email string) (*User, error) {
	var user User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *usersRepository) Update(user *User) error {
	return database.DB.Model(&User{ID: user.ID}).Updates(user).Error
}

func (r *usersRepository) Delete(id uint) error {
	return database.DB.Delete(&User{}, id).Error
}
