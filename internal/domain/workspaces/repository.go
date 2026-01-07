package workspaces

import (
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(workspace *Workspace) error
	FindAll() ([]Workspace, error)
	FindByID(id uint) (*Workspace, error)
	FindByUserID(userID uint) ([]Workspace, error)
	FindByIDs(ids []uint) ([]Workspace, error)
	FindByUserAccess(userID uint) ([]Workspace, error)
	FindGuestWorkspaces(userID uint) ([]Workspace, error)
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
	err := database.DB.
		Select(
			"id",
			"created_by",
			"pass_code",
			"name",
			"privacy",
			"join_link",
			"created_at",
			"updated_at",
		).
		Where("id = ?", id).
		First(&workspace).Error
	return &workspace, err
}

func (r *repository) FindByUserID(userID uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.
		Select("id", "created_by", "pass_code", "name", "privacy", "join_link", "created_at", "updated_at").
		Where("created_by = ?", userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) FindByIDs(ids []uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.
		Select("id", "created_by", "pass_code", "name", "privacy", "join_link", "created_at", "updated_at").
		Where("id IN ?", ids).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) FindByUserAccess(userID uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.Table("workspaces").
		Select("DISTINCT workspaces.*").
		Joins("LEFT JOIN workspaces_users ON workspaces_users.workspace_id = workspaces.id").
		Where("workspaces.created_by = ? OR workspaces_users.user_id = ?", userID, userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) FindGuestWorkspaces(userID uint) ([]Workspace, error) {
	var workspaces []Workspace
	err := database.DB.Table("workspaces").
		Select("DISTINCT workspaces.*").
		Joins("INNER JOIN workspaces_users ON workspaces_users.workspace_id = workspaces.id").
		Where("workspaces_users.user_id = ? AND workspaces.created_by != ?", userID, userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *repository) Update(workspace *Workspace) error {
	return database.DB.Save(workspace).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&Workspace{}, id).Error
}
