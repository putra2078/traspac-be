package user

import (
	"errors"
	"time"

	"hrm-app/internal/domain/contact"
	"hrm-app/internal/pkg/utils"
)

type UseCase interface {
	Register(req *RegisterRequest) error
	GetAll() ([]User, error)
	GetByID(id uint) (*User, error)
	DeleteByID(id uint) error
}

type usecase struct {
	repo        Repository
	contactRepo contact.Repository
}

func NewUseCase(repo Repository, contactRepo contact.Repository) UseCase {
	return &usecase{
		repo:        repo,
		contactRepo: contactRepo,
	}
}

func (u *usecase) Register(req *RegisterRequest) error {
	// 1. Check if email exists
	existing, _ := u.repo.FindByEmail(req.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("email already in use")
	}

	// 2. Validate BirthDate
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return errors.New("invalid birth date format, expected YYYY-MM-DD")
	}

	// 3. Create User
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

	// 4. Create Contact
	newContact := &contact.Contact{
		Name:        req.Name,
		Photo:       req.Photo,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Address:     req.Address,
		BirthDate:   birthDate,
	}

	if err := u.contactRepo.Create(newContact); err != nil {
		// Note: Ideally we should rollback user creation here, but for now we proceed.
		// In a real app, use database transaction.
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

func (u *usecase) Update(user *User) error {
	if _, err := u.repo.FindByID(user.ID); err != nil {
		return err
	}
	return u.repo.Update(user)
}
func (u *usecase) DeleteByID(id uint) error {
	return u.repo.Delete(id)
}
