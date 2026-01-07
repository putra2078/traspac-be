package taskCard

import (
	"context"
	"errors"
)

type UseCase interface {
	Create(ctx context.Context, taskCard *TaskCard) error
	FindAll(ctx context.Context) ([]TaskCard, error)
	FindByID(ctx context.Context, id uint) (*TaskCard, error)
	FindByTaskTabID(ctx context.Context, taskTabID uint) ([]TaskCard, error)
	FindByTaskTabIDs(ctx context.Context, taskTabIDs []uint) ([]TaskCard, error)
	Update(ctx context.Context, taskCard *TaskCard) error
	Delete(ctx context.Context, id uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) Create(ctx context.Context, taskCard *TaskCard) error {
	if taskCard.Name == "" {
		return errors.New("task card name is required")
	}
	return u.repo.Create(ctx, taskCard)
}

func (u *usecase) FindAll(ctx context.Context) ([]TaskCard, error) {
	return u.repo.FindAll(ctx)
}

func (u *usecase) FindByID(ctx context.Context, id uint) (*TaskCard, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *usecase) FindByTaskTabID(ctx context.Context, taskTabID uint) ([]TaskCard, error) {
	return u.repo.FindSummaryByTaskTabIDs(ctx, []uint{taskTabID})
}

func (u *usecase) FindByTaskTabIDs(ctx context.Context, taskTabIDs []uint) ([]TaskCard, error) {
	return u.repo.FindSummaryByTaskTabIDs(ctx, taskTabIDs)
}

func (u *usecase) Update(ctx context.Context, taskCard *TaskCard) error {
	// Check if taskCard exists
	if _, err := u.repo.FindByID(ctx, taskCard.ID); err != nil {
		return err
	}
	return u.repo.Update(ctx, taskCard)
}

func (u *usecase) Delete(ctx context.Context, id uint) error {
	return u.repo.Delete(ctx, id)
}
