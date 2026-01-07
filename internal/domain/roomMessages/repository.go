package room_messages

import (
	"context"
	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(ctx context.Context, message RoomMessage) (RoomMessage, error)
	FindByRoomID(ctx context.Context, roomID uint) ([]RoomMessage, error)
	FindByID(ctx context.Context, id uint) (*RoomMessage, error)
	Update(ctx context.Context, message RoomMessage) (RoomMessage, error)
	Delete(ctx context.Context, id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, message RoomMessage) (RoomMessage, error) {
	if err := database.DB.WithContext(ctx).Create(&message).Error; err != nil {
		return message, err
	}
	// Preload User for the broadcast
	err := database.DB.WithContext(ctx).Preload("User").First(&message, message.ID).Error
	return message, err
}

func (r *repository) FindByRoomID(ctx context.Context, roomID uint) ([]RoomMessage, error) {
	var messages []RoomMessage
	err := database.DB.WithContext(ctx).Preload("User").Where("room_id = ?", roomID).Order("created_at asc").Find(&messages).Error
	return messages, err
}

func (r *repository) FindByID(ctx context.Context, id uint) (*RoomMessage, error) {
	var message RoomMessage
	err := database.DB.WithContext(ctx).Preload("User").First(&message, id).Error
	return &message, err
}

func (r *repository) Update(ctx context.Context, message RoomMessage) (RoomMessage, error) {
	err := database.DB.WithContext(ctx).Save(&message).Error
	if err == nil {
		database.DB.WithContext(ctx).Preload("User").First(&message, message.ID)
	}
	return message, err
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&RoomMessage{}, id).Error
}
