package domain

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	UserID      uint64         `json:"user_id" gorm:"column:id_user;not null;index"`
	Name        string         `json:"name" gorm:"column:nama_toko;not null" validate:"required,min=2,max=100"`
	Description string         `json:"description" gorm:"column:deskripsi"`
	PhotoURL    string         `json:"url_fotol" gorm:"column:url_foto"`
	Status      string         `json:"status" gorm:"default:active"`
	Rating      float64        `json:"rating" gorm:"default:0.0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Store) TableName() string {
	return "toko"
}

type StoreRepository interface {
	Create(store *Store) error
	GetByID(id uint64) (*Store, error)
	GetByUserID(userID uint64) (*Store, error)
	Update(store *Store) error
	Delete(id uint64) error
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
