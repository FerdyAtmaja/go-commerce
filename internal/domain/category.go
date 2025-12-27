package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"column:nama_category;not null;uniqueIndex" validate:"required,min=2,max=100"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CategoryRepository interface {
	Create(category *Category) error
	GetByID(id uint64) (*Category, error)
	GetByName(name string) (*Category, error)
	Update(category *Category) error
	Delete(id uint64) error
	GetAll(limit, offset int) ([]*Category, int64, error)
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}