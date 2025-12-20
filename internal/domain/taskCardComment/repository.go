package taskCardComment

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(taskCardComment *TaskCardComment) error
	FindAll() ([]TaskCardComment, error)
	FindByID(id uint) (*TaskCardComment, error)
	FindByTaskCardID(taskCardID uint) ([]TaskCardComment, error)
	Update(taskCardComment *TaskCardComment) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(taskCardComment *TaskCardComment) error {
	return database.DB.Create(taskCardComment).Error
}

func (r *repository) FindAll() ([]TaskCardComment, error) {
	var taskCardComments []TaskCardComment
	err := database.DB.Find(&taskCardComments).Error
	return taskCardComments, err
}

func (r *repository) FindByID(id uint) (*TaskCardComment, error) {
	var taskCardComment TaskCardComment
	err := database.DB.First(&taskCardComment, id).Error
	return &taskCardComment, err
}

func (r *repository) FindByTaskCardID(taskCardID uint) ([]TaskCardComment, error) {
	var taskCardComments []TaskCardComment
	err := database.DB.Where("task_card_id = ?", taskCardID).Find(&taskCardComments).Error
	return taskCardComments, err
}

func (r *repository) Update(taskCardComment *TaskCardComment) error {
	return database.DB.Model(&TaskCardComment{ID: taskCardComment.ID}).Updates(taskCardComment).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&TaskCardComment{}, id).Error
}