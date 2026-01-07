package taskCardUsers

import (
	"errors"
)

type UseCase interface {
	Create(taskCardUsers *TaskCardUsers) error
	GetByTaskCardID(taskCardID uint) ([]TaskCardUsers, error)
	GetByID(id uint) (*TaskCardUsers, error)
	Update(taskCardUsers *TaskCardUsers) error
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

func (u *usecase) Create(taskCardUsers *TaskCardUsers) error {
	if taskCardUsers == nil {
		return errors.New("payload is required")
	}

	if taskCardUsers.TaskCardID == 0 || taskCardUsers.UserID == 0 {
		return errors.New("task_card_id and user_id are required")
	}

	// Check if user is already assigned to the task card
	existingUser, _ := u.repo.GetByTaskCardIDAndUserID(taskCardUsers.TaskCardID, taskCardUsers.UserID)
	if existingUser != nil && existingUser.ID != 0 {
		return errors.New("user already assigned to this task card")
	}

	return u.repo.Create(taskCardUsers)

}

func (u *usecase) GetByTaskCardID(taskCardID uint) ([]TaskCardUsers, error) {
	return u.repo.GetByTaskCardID(taskCardID)
}

func (u *usecase) GetByID(id uint) (*TaskCardUsers, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) Update(taskCardUsers *TaskCardUsers) error {
	return u.repo.Update(taskCardUsers)
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
