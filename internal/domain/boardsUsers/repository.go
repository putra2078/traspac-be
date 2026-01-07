package boardsUsers

import (
	"hrm-app/internal/pkg/database"

	"gorm.io/gorm"
)

type Repository interface {
	Create(boardUsers *BoardsUsers) error
	GetByBoardID(boardID uint) ([]BoardsUsers, error)
	GetByUserID(userID uint) ([]BoardsUsers, error)
	GetByBoardIDAndUserID(boardID, userID uint) (*BoardsUsers, error)
	GetByID(id uint) (BoardsUsers, error)
	Delete(id uint) error
	Update(boardUsers *BoardsUsers) error
	HasAccess(boardID, userID uint) (bool, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(boardUsers *BoardsUsers) error {
	return database.DB.Create(boardUsers).Error
}

func (r *repository) GetByBoardID(boardID uint) ([]BoardsUsers, error) {
	var boardUsers []BoardsUsers
	err := database.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username")
	}).Where("board_id = ?", boardID).Find(&boardUsers).Error
	return boardUsers, err
}

func (r *repository) GetByUserID(userID uint) ([]BoardsUsers, error) {
	var boardUsers []BoardsUsers
	err := database.DB.Where("user_id = ?", userID).Find(&boardUsers).Error
	return boardUsers, err
}

func (r *repository) GetByBoardIDAndUserID(boardID, userID uint) (*BoardsUsers, error) {
	var boardUser BoardsUsers
	err := database.DB.Preload("User").Where("board_id = ? AND user_id = ?", boardID, userID).First(&boardUser).Error
	return &boardUser, err
}

func (r *repository) GetByID(id uint) (BoardsUsers, error) {
	var boardUsers BoardsUsers
	err := database.DB.Preload("User").First(&boardUsers, id).Error
	return boardUsers, err
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&BoardsUsers{}, id).Error
}

func (r *repository) Update(boardUsers *BoardsUsers) error {
	return database.DB.Save(boardUsers).Error
}

func (r *repository) HasAccess(boardID, userID uint) (bool, error) {
	var count int64
	result := database.DB.Where("board_id = ? AND user_id = ?", boardID, userID).Model(&BoardsUsers{}).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
