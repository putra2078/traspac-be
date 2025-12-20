package taskCard

import (
	"errors"
	"hrm-app/internal/domain/labels"
	"sync"
)

type UseCase interface {
	Create(taskCard *TaskCard) error
	FindAll() ([]TaskCard, error)
	FindByID(id uint) (*TaskCard, error)
	FindByTaskTabID(taskTabID uint) ([]TaskCard, error)
	Update(taskCard *TaskCard) error
	Delete(id uint) error
}

type usecase struct {
	repo      Repository
	labelRepo labels.Repository
}

func NewUseCase(repo Repository, labelRepo labels.Repository) UseCase {
	return &usecase{
		repo:      repo,
		labelRepo: labelRepo,
	}
}

func (u *usecase) Create(taskCard *TaskCard) error {
	if taskCard.Name == "" {
		return errors.New("task card name is required")
	}
	return u.repo.Create(taskCard)
}

func (u *usecase) FindAll() ([]TaskCard, error) {
	taskCards, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for i := range taskCards {
		labels, err := u.labelRepo.FindByTaskCardID(taskCards[i].ID)
		if err != nil {
			return nil, err
		}
		taskCards[i].Labels = labels
	}

	return taskCards, nil
}

func (u *usecase) FindByID(id uint) (*TaskCard, error) {
	var taskCard *TaskCard
	var labels []labels.TaskCardLabel
	var err1, err2 error
	var wg sync.WaitGroup

	wg.Add(2)

	// Fetch task card
	go func() {
		defer wg.Done()
		taskCard, err1 = u.repo.FindByID(id)
	}()

	// Fetch labels
	go func() {
		defer wg.Done()
		labels, err2 = u.labelRepo.FindByTaskCardID(id)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	taskCard.Labels = labels

	return taskCard, nil
}

func (u *usecase) FindByTaskTabID(taskTabID uint) ([]TaskCard, error) {
	taskCards, err := u.repo.FindByTaskTabID(taskTabID)
	if err != nil {
		return nil, err
	}

	for i := range taskCards {
		labels, err := u.labelRepo.FindByTaskCardID(taskCards[i].ID)
		if err != nil {
			return nil, err
		}
		taskCards[i].Labels = labels
	}

	return taskCards, nil
}

func (u *usecase) Update(taskCard *TaskCard) error {
	// Check if taskCard exists
	if _, err := u.repo.FindByID(taskCard.ID); err != nil {
		return err
	}
	return u.repo.Update(taskCard)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
