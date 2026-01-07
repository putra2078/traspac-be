package taskTab

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(taskTab *TaskTab) error
	FindAll() ([]TaskTab, error)
	FindByID(id uint) (*TaskTab, error)
	FindByBoardID(boardID uint) ([]TaskTab, error)
	CreateBatch(taskTabs []TaskTab) error
	Update(taskTab *TaskTab) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(taskTab *TaskTab) error {
	return database.DB.Create(taskTab).Error
}

func (r *repository) FindAll() ([]TaskTab, error) {
	var taskTabs []TaskTab
	err := database.DB.Find(&taskTabs).Error
	return taskTabs, err
}

func (r *repository) FindByID(id uint) (*TaskTab, error) {
	var taskTab TaskTab
	err := database.DB.First(&taskTab, id).Error
	return &taskTab, err
}

func (r *repository) FindByBoardID(boardID uint) ([]TaskTab, error) {
	var taskTabs []TaskTab
	err := database.DB.Where("board_id = ?", boardID).Find(&taskTabs).Error
	return taskTabs, err
}

func (r *repository) FindSummaryByBoardID(boardID uint) ([]TaskTab, error) {
	var taskTabs []TaskTab
	err := database.DB.Select("id, board_id, position, name").Where("board_id = ?", boardID).Find(&taskTabs).Error
	return taskTabs, err
}

func (r *repository) CreateBatch(taskTabs []TaskTab) error {
	return database.DB.Create(taskTabs).Error
}

func (r *repository) Update(taskTab *TaskTab) error {
	return database.DB.Model(&TaskTab{ID: taskTab.ID}).Updates(taskTab).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&TaskTab{}, id).Error
}
