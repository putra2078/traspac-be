package user

import (
	"context"
	"errors"
	"fmt"

	"hrm-app/internal/domain/contact"
	"hrm-app/internal/domain/storage"
	"hrm-app/internal/pkg/utils"
)

type UseCase interface {
	Register(ctx context.Context, req *RegisterRequest) error
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	DeleteByID(id uint) error
	Update(ctx context.Context, id uint, req *UpdateRequest) error
}

type usecase struct {
	repo          Repository
	contactRepo   contact.Repository
	uploadService storage.Service
}

func NewUseCase(repo Repository, contactRepo contact.Repository, uploadService storage.Service) UseCase {
	return &usecase{
		repo:          repo,
		contactRepo:   contactRepo,
		uploadService: uploadService,
	}
}

func (u *usecase) Register(ctx context.Context, req *RegisterRequest) error {
	// 1. Check if email exists
	existing, _ := u.repo.FindByEmail(req.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("email already in use")
	}

	// 2. Create User (Hash Password)
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
	}

	if err := u.repo.Create(newUser); err != nil {
		return err
	}

	// 2.5 Handle Photo Upload
	photoURL := req.Photo // Default to existing URL/string if provided
	if req.PhotoFile != nil {
		bucket := "user-uploads"
		folder := fmt.Sprintf("avatars/users/%d", newUser.ID)
		url, err := u.uploadService.UploadImage(ctx, req.PhotoFile, bucket, folder)
		if err != nil {
			return err
		}
		photoURL = url
	}

	// 3. Create Contact
	newContact := &contact.Contact{
		Name:   req.Name,
		Photo:  photoURL,
		Email:  req.Email,
		UserID: newUser.ID, // Link to User!
	}

	if err := u.contactRepo.Create(ctx, newContact); err != nil {
		// Note: Ideally we should rollback user creation here
		return err
	}

	return nil
}

func (u *usecase) GetAll() ([]User, error) {
	return u.repo.FindAll()
}

func (u *usecase) GetByID(id uint) (*User, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) Update(ctx context.Context, id uint, req *UpdateRequest) error {
	// 1. Get existing user
	user, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2. Get existing contact
	contact, err := u.contactRepo.FindByEmail(ctx, user.Email) // Assumption: Email links them
	// If FindByEmail fails or logic is different (e.g. by UserID if supported), handle it.
	// We added UserID to Contact, so ideally we find by UserID.
	// But repository currently only has FindByEmail. Let's use that for now or assume 1:1 email.
	// Ideally we should add FindByUserID to contact repo, but sticking to existing for minimal changes unless critical.
	// Actually, FindByEmail works if email doesn't change or if we find by old email.
	// Let's assume FindByEmail is safe for now as Register enforces it.
	if err != nil {
		return fmt.Errorf("failed to find contact: %w", err)
	}

	// 3. Update User fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	// Email update usually requires verification, skipping for simplicity or update if needed
	// if req.Email != "" { user.Email = req.Email }

	if err := u.repo.Update(user); err != nil {
		return err
	}

	// 4. Handle Photo Upload
	photoURL := contact.Photo
	if req.Photo != "" {
		photoURL = req.Photo
	} // Direct URL update

	if req.PhotoFile != nil {
		bucket := "user-uploads"
		folder := fmt.Sprintf("avatars/users/%d", user.ID)
		url, err := u.uploadService.UploadImage(ctx, req.PhotoFile, bucket, folder)
		if err != nil {
			return err
		}
		photoURL = url
	}

	// 5. Update Contact fields
	if req.Name != "" {
		contact.Name = req.Name
	}
	contact.Photo = photoURL

	if err := u.contactRepo.Update(ctx, contact); err != nil {
		return err
	}

	return nil
}
func (u *usecase) DeleteByID(id uint) error {
	return u.repo.Delete(id)
}
