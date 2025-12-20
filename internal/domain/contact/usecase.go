package contact

import (
	"errors"
)

type UseCase interface {
	Register(contact *Contact) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) Register(contact *Contact) error {
	existing, _ := u.repo.FindByEmail(contact.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("Email already in use")
	}

	return u.repo.Create(contact)
}
