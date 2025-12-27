package domain

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Description string         `json:"description"`
	PhotoURL    string         `json:"photo_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type StoreRepository interface {
	Create(store *Store) error
	GetByID(id uint) (*Store, error)
	GetByUserID(userID uint) (*Store, error)
	Update(store *Store) error
	Delete(id uint) error
	GetAll(limit, offset int, search string) ([]*Store, int64, error)
}

type CreateStoreRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateStoreRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}