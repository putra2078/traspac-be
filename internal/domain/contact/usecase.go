package contact

import (
	"context"
	"errors"
	"strings"
)

type StorageRepository interface {
	GetURL(bucket string, key string) string
}

type UseCase interface {
	Register(ctx context.Context, contact *Contact) error
	GetByUserID(ctx context.Context, userID uint, bucket string) (*Contact, error)
	Update(ctx context.Context, contact *Contact) error
}

type usecase struct {
	repo        Repository
	storageRepo StorageRepository
}

func NewUseCase(repo Repository, storageRepo StorageRepository) UseCase {
	return &usecase{
		repo:        repo,
		storageRepo: storageRepo,
	}
}

func (u *usecase) Register(ctx context.Context, contact *Contact) error {
	existing, _ := u.repo.FindByEmail(ctx, contact.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("Email already in use")
	}

	return u.repo.Create(ctx, contact)
}

func (u *usecase) GetByUserID(ctx context.Context, userID uint, bucket string) (*Contact, error) {
	contact, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if contact == nil {
		return nil, nil
	}

	// Convert photo path to public URL if it's a storage key
	if contact.Photo != "" && !strings.HasPrefix(contact.Photo, "http") {
		contact.Photo = u.storageRepo.GetURL(bucket, contact.Photo)
	}

	return contact, nil
}

func (u *usecase) Update(ctx context.Context, contact *Contact) error {
	// Check if contact exists
	existing, err := u.repo.FindByUserID(ctx, contact.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("contact not found")
	}

	// Update the contact
	contact.ID = existing.ID
	return u.repo.Update(ctx, contact)
}
