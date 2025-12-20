package boards

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(boards *Boards) error
	FindAll() ([]Boards, error)
	FindByID(id uint) (*Boards, error)
	FindByWorkspaceID(workspaceID uint) ([]Boards, error)
	Update(boards *Boards) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(boards *Boards) error {
	return database.DB.Create(boards).Error
}

func (r *repository) FindAll() ([]Boards, error) {
	var boards []Boards
	err := database.DB.Find(&boards).Error
	return boards, err
}

func (r *repository) FindByID(id uint) (*Boards, error) {
	var boards Boards
	err := database.DB.First(&boards, id).Error
	return &boards, err
}

func (r *repository) FindByWorkspaceID(workspaceID uint) ([]Boards, error) {
	var boards []Boards
	err := database.DB.Where("workspace_id = ?", workspaceID).Find(&boards).Error
	return boards, err
}

func (r *repository) Update(boards *Boards) error {
	return database.DB.Save(boards).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&Boards{}, id).Error
}
