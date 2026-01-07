package contact

import (
	"context"
	"errors"

	"hrm-app/internal/pkg/database"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	FindByEmail(ctx context.Context, email string) (*Contact, error)
	FindByUserID(ctx context.Context, userID uint) (*Contact, error)
	Update(ctx context.Context, contact *Contact) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(ctx context.Context, contact *Contact) error {
	return database.DB.WithContext(ctx).Create(contact).Error
}

func (r *repository) Update(ctx context.Context, contact *Contact) error {
	return database.DB.WithContext(ctx).Save(contact).Error
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*Contact, error) {
	var contact Contact
	err := database.DB.WithContext(ctx).Where("email = ?", email).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jangan return error, biar handler bisa bedain antara "tidak ada data" dan "DB error"
			return &Contact{}, nil
		}
		return nil, err
	}

	return &contact, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID uint) (*Contact, error) {
	var contact Contact
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}
