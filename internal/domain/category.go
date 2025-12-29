package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID                uint64         `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name" gorm:"column:nama_category;not null;uniqueIndex" validate:"required,min=2,max=100"`
	ParentID          *uint64        `json:"parent_id" gorm:"column:parent_id"`
	Slug              string         `json:"slug" gorm:"uniqueIndex;not null"`
	Status            string         `json:"status" gorm:"default:active;check:status IN ('active','inactive')"`
	IsLeaf            bool           `json:"is_leaf" gorm:"default:true"`
	HasChild          bool           `json:"has_child" gorm:"default:false"`
	HasActiveProduct  bool           `json:"has_active_product" gorm:"default:false"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children []*Category `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
}

func (Category) TableName() string {
	return "categories"
}

type CategoryRepository interface {
	Create(category *Category) error
	GetByID(id uint64) (*Category, error)
	GetByName(name string) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	Update(category *Category) error
	Delete(id uint64) error
	GetAll(limit, offset int) ([]*Category, int64, error)
	GetRootCategories(limit, offset int) ([]*Category, int64, error)
	GetChildrenByParentID(parentID uint64) ([]*Category, error)
	HasActiveChildren(categoryID uint64) (bool, error)
	HasActiveProducts(categoryID uint64) (bool, error)
	HasHistoricalProducts(categoryID uint64) (bool, error)
	UpdateStatus(categoryID uint64, status string) error
	GetParentStatus(categoryID uint64) (string, error)
	UpdateHasActiveProduct(categoryID uint64) error
	UpdateChildFlags(categoryID uint64) error
}

type CreateCategoryRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=100" example:"Electronics"`
	ParentID *uint64 `json:"parent_id,omitempty" swaggertype:"integer" example:"0"`
}

type UpdateCategoryRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=100" example:"Updated Electronics"`
	ParentID *uint64 `json:"parent_id,omitempty" swaggertype:"integer" example:"1"`
}