package labels

import (
	"errors"
)

type UseCase interface {
	Create(label *TaskCardLabel) error
	FindAll() ([]TaskCardLabel, error)
	FindByID(id uint) (*TaskCardLabel, error)
	FindByTaskCardID(taskCardID uint) ([]TaskCardLabel, error)
	Update(label *TaskCardLabel) error
	Delete(id uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(label *TaskCardLabel) error {
	if label.Title == "" {
		return errors.New("label title is required")
	}
	return u.repo.Create(label)
}

func (u *usecase) FindAll() ([]TaskCardLabel, error) {
	return u.repo.FindAll()
}

func (u *usecase) FindByID(id uint) (*TaskCardLabel, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) FindByTaskCardID(taskCardID uint) ([]TaskCardLabel, error) {
	return u.repo.FindByTaskCardID(taskCardID)
}

func (u *usecase) Update(label *TaskCardLabel) error {
	// Check if label exists
	if _, err := u.repo.FindByID(label.ID); err != nil {
		return err
	}
	return u.repo.Update(label)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
