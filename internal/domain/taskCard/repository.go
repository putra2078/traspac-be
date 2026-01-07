package taskCard

import (
	"context"
	"hrm-app/internal/pkg/database"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, taskCard *TaskCard) error
	FindAll(ctx context.Context) ([]TaskCard, error)
	FindByID(ctx context.Context, id uint) (*TaskCard, error)
	FindByTaskTabID(ctx context.Context, taskTabID uint) ([]TaskCard, error)
	FindSummaryByTaskTabIDs(ctx context.Context, taskTabIDs []uint) ([]TaskCard, error)
	FindByTaskTabIDPaginated(ctx context.Context, taskTabID uint, limit, offset int) ([]TaskCard, error)
	Update(ctx context.Context, taskCard *TaskCard) error
	Delete(ctx context.Context, id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, taskCard *TaskCard) error {
	return database.DB.WithContext(ctx).Create(taskCard).Error
}

func (r *repository) FindAll(ctx context.Context) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.WithContext(ctx).
		Preload("Labels").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Preload("Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) FindByID(ctx context.Context, id uint) (*TaskCard, error) {
	var taskCard TaskCard
	err := database.DB.WithContext(ctx).
		Preload("Labels").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Preload("Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		First(&taskCard, id).Error
	return &taskCard, err
}

func (r *repository) FindByTaskTabID(ctx context.Context, taskTabID uint) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.WithContext(ctx).
		Preload("Labels").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Preload("Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Where("task_tab_id = ?", taskTabID).
		Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) FindSummaryByTaskTabIDs(ctx context.Context, taskTabIDs []uint) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.WithContext(ctx).
		Preload("Labels").
		Preload("Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Select("id, task_tab_id, name, date, status").
		Where("task_tab_id IN ?", taskTabIDs).
		Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) FindByTaskTabIDPaginated(ctx context.Context, taskTabID uint, limit, offset int) ([]TaskCard, error) {
	var taskCards []TaskCard
	err := database.DB.WithContext(ctx).
		Preload("Labels").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Preload("Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Where("task_tab_id = ?", taskTabID).
		Limit(limit).
		Offset(offset).
		Find(&taskCards).Error
	return taskCards, err
}

func (r *repository) Update(ctx context.Context, taskCard *TaskCard) error {
	return database.DB.WithContext(ctx).Model(&TaskCard{ID: taskCard.ID}).Updates(taskCard).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&TaskCard{}, id).Error
}
