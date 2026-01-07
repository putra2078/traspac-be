package boards

import (
	"context"
	"errors"
	"hrm-app/internal/domain/boardsUsers"
	"hrm-app/internal/domain/labels"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskCardUsers"
	"hrm-app/internal/domain/taskTab"
	"hrm-app/internal/pkg/database"

	"gorm.io/gorm"
)

type UseCase interface {
	Create(ctx context.Context, boards *Boards) error
	FindAll(ctx context.Context) ([]Boards, error)
	FindByID(ctx context.Context, id, userID uint) (*Boards, error)
	FindByWorkspaceID(ctx context.Context, workspaceID uint) ([]Boards, error)
	FindByUserID(ctx context.Context, userID uint) ([]Boards, error)
	Update(ctx context.Context, boards *Boards) error
	Delete(ctx context.Context, id uint) error

	// New methods for optimization
	GetTabsByBoardID(ctx context.Context, boardID uint) ([]TaskTabSummary, error)
	GetCardsByTaskTabID(ctx context.Context, taskTabID uint, limit, offset int) ([]TaskCardSummary, error)
}

type TaskTabSummary struct {
	ID       uint   `json:"id"`
	BoardID  uint   `json:"board_id"`
	Position int    `json:"position"`
	Name     string `json:"name"`
}

type TaskCardSummary struct {
	ID        uint                          `json:"id"`
	TaskTabID uint                          `json:"task_tab_id"`
	Name      string                        `json:"name"`
	Date      string                        `json:"date"`
	Status    bool                          `json:"status"`
	Labels    []labels.TaskCardLabel        `json:"labels"`
	Members   []taskCardUsers.TaskCardUsers `json:"members"`
}

type usecase struct {
	repo              Repository
	taskTabRepo       taskTab.Repository
	taskCardRepo      taskCard.Repository
	boardsUsersRepo   boardsUsers.Repository
	labelsRepo        labels.Repository
	taskCardUsersRepo taskCardUsers.Repository
}

func NewUseCase(
	repo Repository,
	taskTabRepo taskTab.Repository,
	taskCardRepo taskCard.Repository,
	boardsUsersRepo boardsUsers.Repository,
	labelsRepo labels.Repository,
	taskCardUsersRepo taskCardUsers.Repository,
) UseCase {
	return &usecase{
		repo:              repo,
		taskTabRepo:       taskTabRepo,
		taskCardRepo:      taskCardRepo,
		boardsUsersRepo:   boardsUsersRepo,
		labelsRepo:        labelsRepo,
		taskCardUsersRepo: taskCardUsersRepo,
	}
}

func (u *usecase) Create(ctx context.Context, boards *Boards) error {
	if boards.Name == "" {
		return errors.New("name is required")
	}

	err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the board
		if err := tx.Create(boards).Error; err != nil {
			return err
		}

		// Create default task tabs
		defaultTabs := []string{"Todo", "In Progress", "Done"}
		var tabs []taskTab.TaskTab
		// Prealloc slice
		tabs = make([]taskTab.TaskTab, 0, len(defaultTabs))

		for i, tabName := range defaultTabs {
			tabs = append(tabs, taskTab.TaskTab{
				BoardID:  boards.ID,
				Name:     tabName,
				Position: i + 1,
			})
		}
		if err := tx.Create(&tabs).Error; err != nil {
			return err
		}

		// Create board user (creator)
		boardUser := &boardsUsers.BoardsUsers{
			BoardID: boards.ID,
			UserID:  boards.CreatedBy,
		}
		if err := tx.Create(boardUser).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (u *usecase) FindAll(ctx context.Context) ([]Boards, error) {
	return u.repo.FindAll(ctx)
}

func (u *usecase) FindByID(ctx context.Context, id, userID uint) (*Boards, error) {
	// 1. Check Authorization first
	authorizedBoards, err := u.repo.FindByUserAccess(ctx, userID)
	if err != nil {
		return nil, err
	}

	isAuthorized := false
	for _, b := range authorizedBoards {
		if b.ID == id {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return nil, errors.New("unauthorized: you do not have access to this board or board not found")
	}

	// 2. Fetch Board with full details using the optimized Repository JOIN/Preload method
	board, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (u *usecase) FindByUserID(ctx context.Context, userID uint) ([]Boards, error) {
	// Optimization: Fetch ONLY board metadata. No tabs, no cards.
	boards, err := u.repo.FindByUserAccess(ctx, userID)
	if err != nil {
		return nil, err
	}
	// Note: boards come with empty TaskTabs and TaskCards, which is desired for list view.
	return boards, nil
}

func (u *usecase) FindByWorkspaceID(ctx context.Context, workspaceID uint) ([]Boards, error) {
	return u.repo.FindByWorkspaceID(ctx, workspaceID)
}

func (u *usecase) Update(ctx context.Context, boards *Boards) error {
	return u.repo.Update(ctx, boards)
}

func (u *usecase) Delete(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}

func (u *usecase) GetTabsByBoardID(ctx context.Context, boardID uint) ([]TaskTabSummary, error) {
	tabs, err := u.taskTabRepo.FindByBoardID(boardID) // taskTabRepo likely needs Context update too if we want full consistency, checking later
	if err != nil {
		return nil, err
	}

	summaries := make([]TaskTabSummary, 0, len(tabs))
	for _, t := range tabs {
		summaries = append(summaries, TaskTabSummary{
			ID:       t.ID,
			BoardID:  t.BoardID,
			Position: t.Position,
			Name:     t.Name,
		})
	}
	return summaries, nil
}

func (u *usecase) GetCardsByTaskTabID(ctx context.Context, taskTabID uint, limit, offset int) ([]TaskCardSummary, error) {
	// Optimization: Paginated fetch
	cards, err := u.taskCardRepo.FindByTaskTabIDPaginated(ctx, taskTabID, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return []TaskCardSummary{}, nil
	}

	// Convert to summary (which currently matches Entity structure essentially, plus labels/members)
	// FindByTaskTabIDPaginated already preloads Labels and Members.User/Comments.User?
	// Let's check repository implementation. Yes it preloads Labels and Members.User.

	summaries := make([]TaskCardSummary, 0, len(cards))
	for _, c := range cards {
		// Map members to TaskCardUsers struct if needed
		// The repository FindByTaskTabIDPaginated returns []TaskCard.
		// TaskCardSummary expects []taskCardUsers.TaskCardUsers
		// Entity TaskCard has `Members []taskCardUsers.TaskCardUsers`? checking entity..
		// Assuming structure compatibility.

		summaries = append(summaries, TaskCardSummary{
			ID:        c.ID,
			TaskTabID: c.TaskTabID,
			Name:      c.Name,
			Date:      c.Date,
			Status:    c.Status,
			Labels:    c.Labels,
			Members:   c.Members,
		})
	}

	return summaries, nil
}
