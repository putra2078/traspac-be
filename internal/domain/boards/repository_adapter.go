package boards

import (
	"context"
	"hrm-app/internal/domain/boardsUsers"
)

// RepositoryAdapter adapts the full Boards repository to the minimal interface needed by boardsUsers
type RepositoryAdapter struct {
	repo Repository
}

func NewRepositoryAdapter(repo Repository) *RepositoryAdapter {
	return &RepositoryAdapter{repo: repo}
}

// FindByID returns minimal board info needed for authorization
func (a *RepositoryAdapter) FindByID(id uint) (*boardsUsers.BoardInfo, error) {
	board, err := a.repo.FindByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &boardsUsers.BoardInfo{
		ID:          board.ID,
		CreatedBy:   board.CreatedBy,
		WorkspaceID: board.WorkspaceID,
	}, nil
}
