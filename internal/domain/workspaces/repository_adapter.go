package workspaces

import (
	"hrm-app/internal/domain/boardsUsers"
	"hrm-app/internal/domain/workspacesUsers"
)

// repositoryAdapter adapts the workspace Repository to the workspacesUsers.WorkspaceRepository interface
type repositoryAdapter struct {
	repo Repository
}

// NewRepositoryAdapter creates an adapter that implements workspacesUsers.WorkspaceRepository
func NewRepositoryAdapter(repo Repository) workspacesUsers.WorkspaceRepository {
	return &repositoryAdapter{repo: repo}
}

func (a *repositoryAdapter) FindByID(id uint) (*workspacesUsers.WorkspaceInfo, error) {
	workspace, err := a.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &workspacesUsers.WorkspaceInfo{
		ID:        workspace.ID,
		CreatedBy: workspace.CreatedBy,
		PassCode:  workspace.PassCode,
	}, nil
}

// boardWorkspaceRepositoryAdapter adapts the workspace Repository to the boardsUsers.WorkspaceRepository interface
type boardWorkspaceRepositoryAdapter struct {
	repo Repository
}

// NewBoardWorkspaceRepositoryAdapter creates an adapter that implements boardsUsers.WorkspaceRepository
func NewBoardWorkspaceRepositoryAdapter(repo Repository) boardsUsers.WorkspaceRepository {
	return &boardWorkspaceRepositoryAdapter{repo: repo}
}

func (a *boardWorkspaceRepositoryAdapter) FindByID(id uint) (*boardsUsers.WorkspaceInfo, error) {
	workspace, err := a.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &boardsUsers.WorkspaceInfo{
		ID:       workspace.ID,
		PassCode: workspace.PassCode,
	}, nil
}
