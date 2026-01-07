package room_chats

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(roomChats RoomsChats) (RoomsChats, error)
	FindByID(id uint) (RoomsChats, error)
	FindByWorkspaceID(workspaceID uint) ([]RoomsChats, error)
	FindAll() ([]RoomsChats, error)
	Update(roomChats RoomsChats) (RoomsChats, error)
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(roomChats RoomsChats) (RoomsChats, error) {
	err := database.DB.Create(&roomChats).Error
	return roomChats, err
}

func (r *repository) FindByID(id uint) (RoomsChats, error) {
	var roomChats RoomsChats
	err := database.DB.First(&roomChats, id).Error
	return roomChats, err
}

func (r *repository) FindByWorkspaceID(workspaceID uint) ([]RoomsChats, error) {
	var roomChats []RoomsChats
	err := database.DB.Where("workspace_id = ?", workspaceID).Find(&roomChats).Error
	return roomChats, err
}

func (r *repository) FindAll() ([]RoomsChats, error) {
	var roomChats []RoomsChats
	err := database.DB.Find(&roomChats).Error
	return roomChats, err
}

func (r *repository) Update(roomChats RoomsChats) (RoomsChats, error) {
	err := database.DB.Save(&roomChats).Error
	return roomChats, err
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&RoomsChats{}, id).Error
}
