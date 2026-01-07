package room_messages

import (
	"context"
	"errors"
)

type UseCase interface {
	SendMessage(ctx context.Context, message RoomMessage) (RoomMessage, error)
	GetChatHistory(ctx context.Context, roomID uint) ([]RoomMessage, error)
	EditMessage(ctx context.Context, userID uint, messageID uint, newText string) (RoomMessage, error)
	DeleteMessage(ctx context.Context, userID uint, messageID uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) SendMessage(ctx context.Context, message RoomMessage) (RoomMessage, error) {
	return u.repo.Create(ctx, message)
}

func (u *usecase) GetChatHistory(ctx context.Context, roomID uint) ([]RoomMessage, error) {
	return u.repo.FindByRoomID(ctx, roomID)
}

func (u *usecase) EditMessage(ctx context.Context, userID uint, messageID uint, newText string) (RoomMessage, error) {
	// 1. Find message
	msg, err := u.repo.FindByID(ctx, messageID)
	if err != nil {
		return RoomMessage{}, err
	}

	// 2. Check ownership
	if msg.UserID == nil || *msg.UserID != userID {
		return RoomMessage{}, errors.New("unauthorized: not your message")
	}

	// 3. Update text
	msg.MessageText = newText

	// 4. Save
	return u.repo.Update(ctx, *msg)
}

func (u *usecase) DeleteMessage(ctx context.Context, userID uint, messageID uint) error {
	// 1. Find message
	msg, err := u.repo.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	// 2. Check ownership
	if msg.UserID == nil || *msg.UserID != userID {
		return errors.New("unauthorized: not your message")
	}

	// 3. Delete
	return u.repo.Delete(ctx, messageID)
}
