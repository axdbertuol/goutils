// repository.go
package repository

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Getter interface {
	GetAll(entities interface{}, req *http.Request) error
	GetByName(entities interface{}, name string, req *http.Request) error
	GetById(entity interface{}, id uint) error
}

type Paginator interface {
	ApplyPagination(req *http.Request) *gorm.DB
}

type Transactional interface {
	Create(entity interface{}) error
	Update(entity interface{}) error
	Delete(entity interface{}, id uint) error
}

type GormRepository struct {
	DB *gorm.DB
	Getter
	Transactional
	Paginator
}

func ApplyPagination(req *http.Request, db *gorm.DB) *gorm.DB {
	if req.URL.RawQuery != "" {
		return db.Scopes(Paginate(req))
	}
	return db
}

func (gr *GormRepository) ApplyPagination(req *http.Request) *gorm.DB {
	if req.URL.RawQuery != "" {
		return gr.DB.Scopes(Paginate(req))
	}
	return gr.DB
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *GormRepository) GetAll(entities interface{}, req *http.Request) error {
	err := r.DB.Scopes(Paginate(req)).Find(entities).Error
	if err != nil {
		return fmt.Errorf("failed to get all entities: %w", err)
	}
	return nil
}

func (r *GormRepository) GetByName(entities interface{}, name string, req *http.Request) error {

	name = strings.ToLower(name)
	if err := r.DB.
		Scopes(Paginate(req)).
		Where("lower(name) LIKE ?", "%"+name+"%").
		Find(entities).Error; err != nil {
		return fmt.Errorf("failed to get entities by name: %w", err)
	}
	return nil
}

func (r *GormRepository) GetById(entity interface{}, id uint) error {
	if err := r.DB.
		First(entity, id).Error; err != nil {
		return fmt.Errorf("failed to get entity by ID: %w", err)
	}
	return nil
}

func (r *GormRepository) Create(entity interface{}) error {
	if err := r.DB.
		Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

func (r *GormRepository) Update(entity interface{}) error {

	if err := r.DB.
		Model(entity).
		Updates(entity).Error; err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

func (r *GormRepository) Delete(entity interface{}, id uint) error {
	if err := r.DB.Delete(entity, id).Error; err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	return nil
}
