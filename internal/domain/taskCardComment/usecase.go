package taskCardComment

import (
	"errors"
)

type UseCase interface {
	Create(taskCardComment *TaskCardComment) error
	FindAll() ([]TaskCardComment, error)
	FindByID(id uint) (*TaskCardComment, error)
	FindByTaskCardID(taskCardID uint) ([]TaskCardComment, error)
	Update(taskCardComment *TaskCardComment) error
	Delete(id uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) Create(taskCardComment *TaskCardComment) error {
	if taskCardComment.Comment == "" {
		return errors.New("comment is required")
	}
	return u.repo.Create(taskCardComment)
}

func (u *usecase) FindAll() ([]TaskCardComment, error) {
	return u.repo.FindAll()
}

func (u *usecase) FindByID(id uint) (*TaskCardComment, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) Update(taskCardComment *TaskCardComment) error {
	// Check if taskCardComment exists
	if _, err := u.repo.FindByID(uint(taskCardComment.ID)); err != nil {
		return err
	}
	return u.repo.Update(taskCardComment)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) FindByTaskCardID(taskCardID uint) ([]TaskCardComment, error) {
	return u.repo.FindByTaskCardID(taskCardID)
}
