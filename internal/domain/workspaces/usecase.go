package workspaces

import (
	"errors"
	"fmt"
	"hrm-app/config"
	"hrm-app/internal/domain/workspacesUsers"
	"hrm-app/internal/pkg/utils"
)

type UseCase interface {
	Create(workspace *Workspace) error
	GetAll() ([]Workspace, error)
	GetByID(id, userID uint) (*Workspace, error)
	GetByUserID(userID uint) ([]Workspace, error)
	GetGuestWorkspaces(userID uint) ([]Workspace, error)
	DeleteByID(id, userID uint) error
	Update(workspace *Workspace) error
}

type usecase struct {
	repo              Repository
	workSpaceUserRepo workspacesUsers.Repository
	cfg               *config.Config
}

func NewUseCase(repo Repository, workSpaceUserRepo workspacesUsers.Repository, cfg *config.Config) UseCase {
	return &usecase{repo: repo, workSpaceUserRepo: workSpaceUserRepo, cfg: cfg}
}

func (u *usecase) Create(workspace *Workspace) error {
	if workspace.Privacy != "public" && workspace.Privacy != "private" && workspace.Privacy != "team" {
		return errors.New("privacy must be either 'public', 'private', or 'team'")
	}

	workspace.PassCode = utils.GeneratePassCode(6)

	if err := u.repo.Create(workspace); err != nil {
		return err
	}

	// Generate Join Token and update JoinLink
	token, err := utils.GenerateJoinToken(u.cfg, workspace.ID, "workspace", workspace.PassCode)
	if err == nil {
		// Example URL, you can adjust the base URL as needed
		workspace.JoinLink = fmt.Sprintf("http://be.putratek.my.id/api/v1/workspaces/join?token=%s", token)
		_ = u.repo.Update(workspace)
	}

	// Add creator to workspace users
	workspaceUsers := &workspacesUsers.WorkspacesUsers{
		WorkspaceID: workspace.ID,
		UserID:      workspace.CreatedBy,
	}

	return u.workSpaceUserRepo.Create(workspaceUsers)
}

func (u *usecase) GetAll() ([]Workspace, error) {
	return u.repo.FindAll()
}

func (u *usecase) GetByID(id, userID uint) (*Workspace, error) {
	workspaces, err := u.repo.FindByUserAccess(userID)
	if err != nil {
		return nil, err
	}

	for _, w := range workspaces {
		if w.ID == id {
			return &w, nil
		}
	}

	return nil, errors.New("unauthorized: you are not a member of this workspace or workspace not found")
}

func (u *usecase) GetByUserID(userID uint) ([]Workspace, error) {
	return u.repo.FindByUserAccess(userID)
}

func (u *usecase) GetGuestWorkspaces(userID uint) ([]Workspace, error) {
	return u.repo.FindGuestWorkspaces(userID)
}

func (u *usecase) DeleteByID(id, userID uint) error {
	workspace, err := u.repo.FindByID(id)
	if err != nil || workspace == nil {
		return errors.New("workspace not found")
	}

	if workspace.CreatedBy != userID {
		return errors.New("unauthorized: only the workspace creator can delete this workspace")
	}

	return u.repo.Delete(id)
}

func (u *usecase) Update(workspace *Workspace) error {
	if workspace.Privacy != "public" && workspace.Privacy != "private" && workspace.Privacy != "team" {
		return errors.New("privacy must be either 'public', 'private', or 'team'")
	}
	return u.repo.Update(workspace)
}
