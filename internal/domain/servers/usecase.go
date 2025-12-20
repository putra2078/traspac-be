package servers

import (
	"errors"
)

type UseCase interface {
	Create(server *Server) error
	GetAll() ([]Server, error)
	GetByID(id uint) (*Server, error)
	Update(server *Server) error
	DeleteByID(id uint) error
}

type usecase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(server *Server) error {
	existing, _ := u.repo.FindByID(server.ID)
	if existing != nil && existing.ID != 0 {
		return errors.New("ID already in use")
	}

	return u.repo.Create(server)
}

func (u *usecase) GetAll() ([]Server, error) {
	return u.repo.FindAll()
}

func (u *usecase) GetByID(id uint) (*Server, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) DeleteByID(id uint) error {
	return u.repo.Delete(id)
}

func (u *usecase) Update(server *Server) error {
	// Check if server exists
	if _, err := u.repo.FindByID(server.ID); err != nil {
		return err
	}
	return u.repo.Update(server)
}
