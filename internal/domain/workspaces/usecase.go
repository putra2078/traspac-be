package workspaces

import (
	"errors"

	"hrm-app/internal/pkg/utils"
)



type UseCase interface {
	Create(workspace *Workspace) error
	GetAll() ([]Workspace, error)
	GetByID(id uint) (*Workspace, error)
	GetByUserID(userID uint) ([]Workspace, error)
	DeleteByID(id uint) error
	Update(workspace *Workspace) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(workspace *Workspace) error {
	if workspace.Privacy != "public" && workspace.Privacy != "private" && workspace.Privacy != "team" {
		return errors.New("privacy must be either 'public', 'private', or 'team'")
	}

	workspace.PassCode = utils.GeneratePassCode(6)
	workspace.JoinLink = utils.GenerateLinkJoin()

	return u.repo.Create(workspace)
}

func (u *usecase) GetAll() ([]Workspace, error) {
	return u.repo.FindAll()
}

func (u *usecase) GetByID(id uint) (*Workspace, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) GetByUserID(userID uint) ([]Workspace, error) {
	return u.repo.FindByUserID(userID)
}

func (u *usecase) DeleteByID(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) Update(workspace *Workspace) error {
	if workspace.Privacy != "public" && workspace.Privacy != "private" && workspace.Privacy != "team" {
		return errors.New("privacy must be either 'public', 'private', or 'team'")
	}
	return u.repo.Update(workspace)
}
