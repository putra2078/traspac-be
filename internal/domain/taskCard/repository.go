package taskCard

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(taskCard *TaskCard) error
	FindAll() ([]TaskCard, error)
	FindByID(id uint) (*TaskCard, error)
	FindByTaskTabID(taskTabID uint) ([]TaskCard, error)
	FindSummaryByTaskTabID(taskTabID uint) ([]TaskCard, error)
	Update(taskCard *TaskCard) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(taskCard *TaskCard) error {
	return database.DB.Create(taskCard).Error
}

func (r *repository) FindAll() ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) FindByID(id uint) (*TaskCard, error) {
	var taskCard TaskCard
	err := database.DB.Preload("Labels").First(&taskCard, id).Error
	return &taskCard, err
}

func (r *repository) FindByTaskTabID(taskTabID uint) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.Preload("Labels").Where("task_tab_id = ?", taskTabID).Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) FindSummaryByTaskTabID(taskTabID uint) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.Preload("Labels").Select("id, task_tab_id, name, date, status").Where("task_tab_id = ?", taskTabID).Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) Update(taskCard *TaskCard) error {
	return database.DB.Model(&TaskCard{ID: taskCard.ID}).Updates(taskCard).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&TaskCard{}, id).Error
}
