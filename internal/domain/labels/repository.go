package labels

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(label *TaskCardLabel) error
	FindAll() ([]TaskCardLabel, error)
	FindByID(id uint) (*TaskCardLabel, error)
	FindByTaskCardID(taskCardID uint) ([]TaskCardLabel, error)
	FindByTaskCardIDs(taskCardIDs []uint) ([]TaskCardLabel, error)
	Update(label *TaskCardLabel) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(label *TaskCardLabel) error {
	return database.DB.Create(label).Error
}

func (r *repository) FindAll() ([]TaskCardLabel, error) {
	var labels []TaskCardLabel
	err := database.DB.Find(&labels).Error
	return labels, err
}

func (r *repository) FindByID(id uint) (*TaskCardLabel, error) {
	var label TaskCardLabel
	err := database.DB.First(&label, id).Error
	return &label, err
}

func (r *repository) FindByTaskCardID(taskCardID uint) ([]TaskCardLabel, error) {
	var labels []TaskCardLabel
	err := database.DB.Where("task_card_id = ?", taskCardID).Find(&labels).Error
	return labels, err
}

func (r *repository) FindByTaskCardIDs(taskCardIDs []uint) ([]TaskCardLabel, error) {
	var labels []TaskCardLabel
	err := database.DB.Where("task_card_id IN ?", taskCardIDs).Find(&labels).Error
	return labels, err
}

func (r *repository) Update(label *TaskCardLabel) error {
	return database.DB.Model(&TaskCardLabel{ID: label.ID}).Updates(label).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&TaskCardLabel{}, id).Error
}
