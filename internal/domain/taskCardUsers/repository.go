package taskCardUsers

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(taskCardUsers *TaskCardUsers) error
	GetByTaskCardID(taskCardID uint) ([]TaskCardUsers, error)
	GetByTaskCardIDs(taskCardIDs []uint) ([]TaskCardUsers, error)
	GetByID(id uint) (*TaskCardUsers, error)
	GetByTaskCardIDAndUserID(taskCardID, userID uint) (*TaskCardUsers, error)
	Update(taskCardUsers *TaskCardUsers) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(taskCardUsers *TaskCardUsers) error {
	return database.DB.Create(taskCardUsers).Error
}

func (r *repository) GetByTaskCardID(taskCardID uint) ([]TaskCardUsers, error) {
	var taskCardUsers []TaskCardUsers
	err := database.DB.Preload("User").Where("task_card_id = ?", taskCardID).Find(&taskCardUsers).Error
	return taskCardUsers, err
}

func (r *repository) GetByTaskCardIDs(taskCardIDs []uint) ([]TaskCardUsers, error) {
	var taskCardUsers []TaskCardUsers
	err := database.DB.Preload("User").Where("task_card_id IN ?", taskCardIDs).Find(&taskCardUsers).Error
	return taskCardUsers, err
}

func (r *repository) GetByID(id uint) (*TaskCardUsers, error) {
	var taskCardUser TaskCardUsers
	err := database.DB.Preload("User").First(&taskCardUser, id).Error
	return &taskCardUser, err
}

func (r *repository) GetByTaskCardIDAndUserID(taskCardID, userID uint) (*TaskCardUsers, error) {
	var taskCardUser TaskCardUsers
	err := database.DB.Preload("User").Where("task_card_id = ? AND user_id = ?", taskCardID, userID).First(&taskCardUser).Error
	return &taskCardUser, err
}

func (r *repository) Update(taskCardUsers *TaskCardUsers) error {
	return database.DB.Model(&TaskCardUsers{ID: taskCardUsers.ID}).Updates(taskCardUsers).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&TaskCardUsers{}, id).Error
}
