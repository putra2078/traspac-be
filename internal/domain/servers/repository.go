package servers

import (
	// "time"

	"hrm-app/internal/pkg/database"
)

type Repository interface {
	Create(server *Server) error
	FindAll() ([]Server, error)
	FindByID(id uint) (*Server, error)
	Update(server *Server) error
	Delete(id uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(server *Server) error {
	return database.DB.Create(server).Error
}

func (r *repository) FindAll() ([]Server, error) {
	var servers []Server
	err := database.DB.Find(&servers).Error
	return servers, err
}

func (r *repository) FindByID(id uint) (*Server, error) {
	var server Server
	err := database.DB.First(&server, id).Error
	return &server, err
}

func (r *repository) Update(server *Server) error {
	return database.DB.Model(&Server{ID: server.ID}).Updates(server).Error
}

func (r *repository) Delete(id uint) error {
	return database.DB.Delete(&Server{}, id).Error
}
