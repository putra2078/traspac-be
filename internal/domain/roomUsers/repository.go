package roomUsers

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(roomUser *RoomUsers) error
	Delete(id uint) error
	FindByRoomID(roomID uint) ([]RoomUsers, error)
	FindByUserID(userID uint) ([]RoomUsers, error)
	FindByRoomIDAndUserID(roomID, userID uint) (*RoomUsers, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(roomUser *RoomUsers) error {
	return database.DB.Create(roomUser).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&RoomUsers{}, id).Error
}

func (r *repository) FindByRoomID(roomID uint) ([]RoomUsers, error) {
	var roomUsers []RoomUsers
	err := database.DB.Preload("User").Where("room_id = ?", roomID).Find(&roomUsers).Error
	return roomUsers, err
}

func (r *repository) FindByUserID(userID uint) ([]RoomUsers, error) {
	var roomUsers []RoomUsers
	err := database.DB.Where("user_id = ?", userID).Find(&roomUsers).Error
	return roomUsers, err
}

func (r *repository) FindByRoomIDAndUserID(roomID, userID uint) (*RoomUsers, error) {
	var roomUser RoomUsers
	err := database.DB.Where("room_id = ? AND user_id = ?", roomID, userID).First(&roomUser).Error
	if err != nil {
		return nil, err
	}
	return &roomUser, nil
}
