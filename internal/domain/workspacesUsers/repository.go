package workspacesUsers

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(workspacesUsers *WorkspacesUsers) error
	GetByWorkspaceID(workspaceID uint) ([]WorkspacesUsers, error)
	GetByUserID(userID uint) ([]WorkspacesUsers, error)
	GetByWorkspaceIDAndUserID(workspaceID, userID uint) (*WorkspacesUsers, error)
	GetByID(id uint) (WorkspacesUsers, error)
	Delete(id uint) error
	Update(workspacesUsers *WorkspacesUsers) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(workspacesUsers *WorkspacesUsers) error {
	return database.DB.Create(workspacesUsers).Error
}

func (r *repository) GetByWorkspaceID(workspaceID uint) ([]WorkspacesUsers, error) {
	var workspacesUsers []WorkspacesUsers
	if err := database.DB.Preload("User").Where("workspace_id = ?", workspaceID).Find(&workspacesUsers).Error; err != nil {
		return nil, err
	}
	return workspacesUsers, nil
}

func (r *repository) GetByUserID(userID uint) ([]WorkspacesUsers, error) {
	var workspacesUsers []WorkspacesUsers
	if err := database.DB.Preload("User").Where("user_id = ?", userID).Find(&workspacesUsers).Error; err != nil {
		return nil, err
	}
	return workspacesUsers, nil
}

func (r *repository) GetByWorkspaceIDAndUserID(workspaceID, userID uint) (*WorkspacesUsers, error) {
	var workspaceUser WorkspacesUsers
	if err := database.DB.Preload("User").Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&workspaceUser).Error; err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}

func (r *repository) GetByID(id uint) (WorkspacesUsers, error) {
	var workspacesUsers WorkspacesUsers
	if err := database.DB.Preload("User").Where("id = ?", id).Find(&workspacesUsers).Error; err != nil {
		return WorkspacesUsers{}, err
	}
	return workspacesUsers, nil
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&WorkspacesUsers{}, id).Error
}

func (r *repository) Update(workspacesUsers *WorkspacesUsers) error {
	return database.DB.Save(workspacesUsers).Error
}
