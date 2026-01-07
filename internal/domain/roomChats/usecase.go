package room_chats

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"hrm-app/internal/domain/storage"
)

type UseCase interface {
	Create(roomChats *RoomsChats) error
	GetByID(id uint) (*RoomsChats, error)
	GetByWorkspaceID(workspaceID uint) ([]RoomsChats, error)
	GetAll() ([]RoomsChats, error)
	Update(roomChats *RoomsChats) error
	Delete(id uint) error
	UploadAttachment(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type usecase struct {
	repo          Repository
	uploadService storage.Service
	bucketName    string
}

func NewUseCase(repo Repository, uploadService storage.Service, bucketName string) UseCase {
	return &usecase{
		repo:          repo,
		uploadService: uploadService,
		bucketName:    bucketName,
	}
}

func (u *usecase) Create(roomChats *RoomsChats) error {
	if roomChats.Name == "" {
		return errors.New("room chat name is required")
	}
	_, err := u.repo.Create(*roomChats)
	return err
}

func (u *usecase) GetByID(id uint) (*RoomsChats, error) {
	data, err := u.repo.FindByID(id)
	return &data, err
}

func (u *usecase) GetByWorkspaceID(workspaceID uint) ([]RoomsChats, error) {
	return u.repo.FindByWorkspaceID(workspaceID)
}

func (u *usecase) GetAll() ([]RoomsChats, error) {
	return u.repo.FindAll()
}

func (u *usecase) Update(roomChats *RoomsChats) error {
	// Check if room chat exists
	if _, err := u.repo.FindByID(roomChats.ID); err != nil {
		return err
	}
	_, err := u.repo.Update(*roomChats)
	return err
}

func (u *usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) UploadAttachment(ctx context.Context, file *multipart.FileHeader) (string, error) {
	bucket := u.bucketName
	folder := "room-chats"

	// Upload using common storage service
	// The service handles generating unique filenames, etc.
	url, err := u.uploadService.UploadImage(ctx, file, bucket, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload attachment: %w", err)
	}
	return url, nil
}
