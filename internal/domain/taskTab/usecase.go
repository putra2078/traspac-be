package taskTab

import (
	"errors"
)

type UseCase interface {
	Create(taskTab *TaskTab) error
	FindAll() ([]TaskTab, error)
	FindByID(id uint) (*TaskTab, error)
	Update(taskTab *TaskTab) error
	Delete(id uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repository Repository) UseCase {
	return &usecase{repo: repository}
}

func (u *usecase) Create(taskTab *TaskTab) error {
	if taskTab.Name == "" {
		return errors.New("name is required")
	}
	return u.repo.Create(taskTab)
}

func (u *usecase) FindAll() ([]TaskTab, error) {
	return u.repo.FindAll()
}

func (u *usecase) FindByID(id uint) (*TaskTab, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) Update(taskTab *TaskTab) error {
	// Check if taskTab exists
	if _, err := u.repo.FindByID(taskTab.ID); err != nil {
		return err
	}
	return u.repo.Update(taskTab)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
