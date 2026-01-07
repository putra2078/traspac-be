package roomUsers

import (
	"errors"
)

type UseCase interface {
	Join(roomID, userID uint) error
	Leave(roomID, userID uint) error
	GetUsersByRoom(roomID uint) ([]RoomUsers, error)
	IsUserInRoom(roomID, userID uint) (bool, error)
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) Join(roomID, userID uint) error {
	// Check if already in room
	exist, _ := u.repo.FindByRoomIDAndUserID(roomID, userID)
	if exist != nil {
		return nil // Already joined
	}

	roomUser := &RoomUsers{
		RoomID: roomID,
		UserID: userID,
	}
	return u.repo.Create(roomUser)
}

func (u *usecase) Leave(roomID, userID uint) error {
	roomUser, err := u.repo.FindByRoomIDAndUserID(roomID, userID)
	if err != nil {
		return errors.New("membership not found")
	}
	return u.repo.Delete(roomUser.ID)
}

func (u *usecase) GetUsersByRoom(roomID uint) ([]RoomUsers, error) {
	return u.repo.FindByRoomID(roomID)
}

func (u *usecase) IsUserInRoom(roomID, userID uint) (bool, error) {
	roomUser, err := u.repo.FindByRoomIDAndUserID(roomID, userID)
	if err != nil {
		return false, nil
	}
	return roomUser != nil, nil
}
