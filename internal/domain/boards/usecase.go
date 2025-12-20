package boards

import (
	"errors"
	"hrm-app/internal/domain/taskCard"
	"hrm-app/internal/domain/taskTab"
	"sync"
)

type UseCase interface {
	Create(boards *Boards) error
	FindAll() ([]Boards, error)
	FindByID(id uint) (*Boards, error)
	FindByWorkspaceID(workspaceID uint) ([]Boards, error)
	Update(boards *Boards) error
	Delete(id uint) error
}

type TaskTabSummary struct {
	ID       uint   `json:"id"`
	BoardID  uint   `json:"board_id"`
	Position int    `json:"position"`
	Name     string `json:"name"`
}

type TaskCardSummary struct {
	ID        uint   `json:"id"`
	TaskTabID uint   `json:"task_tab_id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	Status    bool   `json:"status"`
}

type usecase struct {
	repo         Repository
	taskTabRepo  taskTab.Repository
	taskCardRepo taskCard.Repository
}

func NewUseCase(repo Repository, taskTabRepo taskTab.Repository, taskCardRepo taskCard.Repository) UseCase {
	return &usecase{
		repo:         repo,
		taskTabRepo:  taskTabRepo,
		taskCardRepo: taskCardRepo,
	}
}

func (u *usecase) Create(boards *Boards) error {
	if boards.Name == "" {
		return errors.New("name is required")
	}

	// Create the board first
	if err := u.repo.Create(boards); err != nil {
		return err
	}

	// Create default task tabs using goroutines
	defaultTabs := []string{"Todo", "In Progress", "Done"}
	var wg sync.WaitGroup
	errChan := make(chan error, len(defaultTabs))

	for i, tabName := range defaultTabs {
		wg.Add(1)
		go func(name string, position int) {
			defer wg.Done()
			if err := u.taskTabRepo.Create(&taskTab.TaskTab{
				BoardID:  boards.ID,
				Name:     name,
				Position: position,
			}); err != nil {
				errChan <- err
			}
		}(tabName, i+1)
	}

	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *usecase) FindAll() ([]Boards, error) {
	return u.repo.FindAll()
}

func (u *usecase) FindByID(id uint) (*Boards, error) {
	board, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Fetch task tabs first
	// Fetch task tabs first
	tabs, err := u.taskTabRepo.FindSummaryByBoardID(board.ID)
	if err != nil {
		return nil, err
	}

	var tabSummaries []TaskTabSummary
	for _, t := range tabs {
		tabSummaries = append(tabSummaries, TaskTabSummary{
			ID:       t.ID,
			BoardID:  t.BoardID,
			Position: t.Position,
			Name:     t.Name,
		})
	}
	board.TaskTabs = tabSummaries

	// Fetch task cards for each tab concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allCards []TaskCardSummary

	wg.Add(len(tabs))
	for _, tab := range tabs {
		go func(tabID uint) {
			defer wg.Done()
			cards, err := u.taskCardRepo.FindSummaryByTaskTabID(tabID)
			if err == nil {
				var cardSummaries []TaskCardSummary
				for _, c := range cards {
					cardSummaries = append(cardSummaries, TaskCardSummary{
						ID:        c.ID,
						TaskTabID: c.TaskTabID,
						Name:      c.Name,
						Date:      c.Date,
						Status:    c.Status,
					})
				}
				mu.Lock()
				allCards = append(allCards, cardSummaries...)
				mu.Unlock()
			}
		}(tab.ID)
	}

	wg.Wait()
	board.TaskCards = allCards
	return board, nil
}

func (u *usecase) FindByWorkspaceID(workspaceID uint) ([]Boards, error) {
	return u.repo.FindByWorkspaceID(workspaceID)
}

func (u *usecase) Update(boards *Boards) error {
	return u.repo.Update(boards)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
