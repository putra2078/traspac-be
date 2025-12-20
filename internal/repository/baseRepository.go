package repository

import (
	"errors"

	"gorm.io/gorm"
)

type Pagination struct {
	Page       int
	Limit      int
	TotalRows  int64
	TotalPages int
}

type BaseRepository[T any] interface {
	Create(entity *T) error
	FindAll(entities *[]T) error
	FindByID(id uint, entity *T) error
	Update(entity *T) error
	Delete(id uint) error
	FindWithPagination(entities *[]T, page int, limit int, pagination *Pagination) error
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *baseRepository[T]) FindAll(entities *[]T) error {
	return r.db.Find(entities).Error
}

func (r *baseRepository[T]) FindByID(id uint, entity *T) error {
	err := r.db.First(entity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("data not found")
	}
	return err
}

func (r *baseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *baseRepository[T]) Delete(id uint) error {
	// Soft delete default dari GORM akan aktif jika model punya gorm.DeletedAt
	return r.db.Delete(new(T), id).Error
}

func (r *baseRepository[T]) FindWithPagination(entities *[]T, page int, limit int, pagination *Pagination) error {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Count total rows
	r.db.Model(new(T)).Count(&pagination.TotalRows)

	// Query data
	err := r.db.Limit(limit).Offset(offset).Find(entities).Error
	if err != nil {
		return err
	}

	pagination.Page = page
	pagination.Limit = limit

	// Hitung total pages
	if pagination.TotalRows%int64(limit) == 0 {
		pagination.TotalPages = int(pagination.TotalRows / int64(limit))
	} else {
		pagination.TotalPages = int((pagination.TotalRows / int64(limit)) + 1)
	}

	return nil
}
