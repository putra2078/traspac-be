package workspacesUsers

import (
	"errors"
	"hrm-app/config"
	"hrm-app/internal/pkg/utils"
)

// WorkspaceRepository defines minimal interface needed to verify workspace ownership
type WorkspaceRepository interface {
	FindByID(id uint) (*WorkspaceInfo, error)
}

// WorkspaceInfo contains minimal workspace information needed for authorization
type WorkspaceInfo struct {
	ID        uint
	CreatedBy uint
	PassCode  string
}

type UseCase interface {
	Create(workspacesUsers *WorkspacesUsers, requestingUserID uint) error
	GetByWorkspaceID(workspaceID uint) ([]WorkspacesUsers, error)
	GetByUserID(userID uint) ([]WorkspacesUsers, error)
	GetByID(id uint) (WorkspacesUsers, error)
	Delete(id uint) error
	Update(workspacesUsers *WorkspacesUsers) error
	Join(userID uint, token string) error
	GenerateJoinToken(workspaceID, userID uint) (string, error)
}

type usecase struct {
	repo          Repository
	workspaceRepo WorkspaceRepository
	cfg           *config.Config
}

func NewUseCase(repo Repository, workspaceRepo WorkspaceRepository, cfg *config.Config) UseCase {
	return &usecase{
		repo:          repo,
		workspaceRepo: workspaceRepo,
		cfg:           cfg,
	}
}

func (u *usecase) Create(workspacesUsers *WorkspacesUsers, requestingUserID uint) error {
	// Fetch the workspace to verify ownership
	workspace, err := u.workspaceRepo.FindByID(workspacesUsers.WorkspaceID)

	if err != nil {
		return errors.New("workspace not found")
	}

	// Check if the requesting user is the creator of the workspace
	if workspace.CreatedBy != requestingUserID {
		return errors.New("unauthorized: only workspace creator can add users")
	}

	// Check if user is already assigned to the workspace
	existingUser, _ := u.repo.GetByWorkspaceIDAndUserID(workspacesUsers.WorkspaceID, workspacesUsers.UserID)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already assigned to this workspace")
	}

	return u.repo.Create(workspacesUsers)
}

func (u *usecase) GetByWorkspaceID(workspaceID uint) ([]WorkspacesUsers, error) {
	return u.repo.GetByWorkspaceID(workspaceID)
}

func (u *usecase) GetByUserID(userID uint) ([]WorkspacesUsers, error) {
	return u.repo.GetByUserID(userID)
}

func (u *usecase) GetByID(id uint) (WorkspacesUsers, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) Update(workspacesUsers *WorkspacesUsers) error {
	return u.repo.Update(workspacesUsers)
}

func (u *usecase) GenerateJoinToken(workspaceID, userID uint) (string, error) {
	workspace, err := u.workspaceRepo.FindByID(workspaceID)
	if err != nil {
		return "", errors.New("workspace not found")
	}

	// Only workspace creator can generate join token
	if workspace.CreatedBy != userID {
		return "", errors.New("unauthorized: only workspace creator can generate join token")
	}

	return utils.GenerateJoinToken(u.cfg, workspaceID, "workspace", workspace.PassCode)
}

func (u *usecase) Join(userID uint, token string) error {
	claims, err := utils.ValidateJoinToken(u.cfg, token)
	if err != nil {
		return err
	}

	if claims.EntityType != "workspace" {
		return errors.New("invalid token type")
	}

	workspaceID := claims.EntityID
	passCode := claims.PassCode

	// Fetch the workspace to verify passcode
	workspace, err := u.workspaceRepo.FindByID(workspaceID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Check if the passcode matches
	if workspace.PassCode != passCode {
		return errors.New("invalid passcode")
	}

	// Check if user is already assigned to the workspace
	existingUser, _ := u.repo.GetByWorkspaceIDAndUserID(workspaceID, userID)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already assigned to this workspace")
	}

	workspacesUsers := &WorkspacesUsers{
		WorkspaceID: workspaceID,
		UserID:      userID,
	}

	return u.repo.Create(workspacesUsers)
}
