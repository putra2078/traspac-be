package boards

import (
	"context"
	"hrm-app/internal/pkg/database"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, boards *Boards) error
	FindAll(ctx context.Context) ([]Boards, error)
	FindByID(ctx context.Context, id uint) (*Boards, error)
	FindByWorkspaceID(ctx context.Context, workspaceID uint) ([]Boards, error)
	FindByUserID(ctx context.Context, userID uint) ([]Boards, error)
	FindByIDs(ctx context.Context, ids []uint) ([]Boards, error)
	FindByUserAccess(ctx context.Context, userID uint) ([]Boards, error)
	Update(ctx context.Context, boards *Boards) error
	Delete(ctx context.Context, id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, boards *Boards) error {
	return database.DB.WithContext(ctx).Create(boards).Error
}

func (r *repository) FindAll(ctx context.Context) ([]Boards, error) {
	var boards []Boards
	err := database.DB.WithContext(ctx).Find(&boards).Error
	return boards, err
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Boards, error) {
	var boards Boards
	err := database.DB.WithContext(ctx).
		Preload("TaskTabs", func(db *gorm.DB) *gorm.DB {
			return db.Order("position asc")
		}).
		Preload("TaskTabs.TaskCards").
		Preload("TaskTabs.TaskCards.Labels").
		Preload("TaskTabs.TaskCards.Members.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username")
		}).
		Select("id", "workspace_id", "created_by", "name", "images", "created_at", "updated_at").
		Where("id = ?", id).
		First(&boards).Error

	if err == nil {
		// Populate flat TaskCards list from nested structure
		// This ensures compatibility with frontend expecting board.task_cards
		var allCards []interface{}
		for _, tab := range boards.TaskTabs {
			for _, card := range tab.TaskCards {
				allCards = append(allCards, card)
			}
		}
		boards.TaskCards = allCards
	}

	return &boards, err
}

func (r *repository) FindByWorkspaceID(ctx context.Context, workspaceID uint) ([]Boards, error) {
	var boards []Boards
	err := database.DB.WithContext(ctx).
		Select("id", "workspace_id", "created_by", "name", "images").
		Where("workspace_id = ?", workspaceID).
		Find(&boards).Error
	return boards, err
}

func (r *repository) FindByUserID(ctx context.Context, userID uint) ([]Boards, error) {
	var boards []Boards
	err := database.DB.WithContext(ctx).
		Select("id", "workspace_id", "created_by", "name", "images").
		Where("created_by = ?", userID).
		Find(&boards).Error
	return boards, err
}

func (r *repository) FindByIDs(ctx context.Context, ids []uint) ([]Boards, error) {
	var boards []Boards
	err := database.DB.WithContext(ctx).
		Select("id", "workspace_id", "created_by", "name", "images", "created_at", "updated_at").
		Where("id IN ?", ids).
		Find(&boards).Error
	return boards, err
}

func (r *repository) FindByUserAccess(ctx context.Context, userID uint) ([]Boards, error) {
	var results []Boards
	err := database.DB.WithContext(ctx).Table("boards").
		Select("DISTINCT boards.*").
		Joins("LEFT JOIN boards_users ON boards_users.board_id = boards.id").
		Where("boards.created_by = ? OR boards_users.user_id = ?", userID, userID).
		Find(&results).Error
	return results, err
}

func (r *repository) Update(ctx context.Context, boards *Boards) error {
	return database.DB.WithContext(ctx).Save(boards).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&Boards{}, id).Error
}
