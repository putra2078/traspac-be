package workspaces

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(workspace *Workspace) error
	FindAll() ([]Workspace, error)
	FindByID(id uint) (*Workspace, error)
	FindByUserID(userID uint) ([]Workspace, error)
	Update(workspace *Workspace) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(workspace *Workspace) error {
	return database.DB.Create(workspace).Error
}

func (r *repository) FindAll() ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) FindByID(id uint) (*Workspace, error) {
	var workspace Workspace
	err := database.DB.First(&workspace, id).Error
	return &workspace, err
}

func (r *repository) FindByUserID(userID uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.Where("created_by = ?", userID).Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) Update(workspace *Workspace) error {
	return database.DB.Save(workspace).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&Workspace{}, id).Error
}