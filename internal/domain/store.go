package domain

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID          uint64         `json:"id" gorm:"primaryKey;column:id"`
	UserID      uint64         `json:"user_id" gorm:"column:id_user;type:bigint unsigned;not null;uniqueIndex:idx_toko_user"`
	Name        string         `json:"name" gorm:"column:nama_toko;type:varchar(255);not null" validate:"required,min=2,max=255"`
	PhotoURL    string         `json:"url_foto" gorm:"column:url_foto;type:varchar(255)"`
	Description string         `json:"description" gorm:"column:deskripsi;type:text"`
	Status      string         `json:"status" gorm:"column:status;type:enum('pending','active','inactive','suspended');default:pending;index:idx_toko_status" validate:"omitempty,oneof=pending active inactive suspended"`
	Rating      float64        `json:"rating" gorm:"column:rating;type:decimal(2,1);default:0.0;index:idx_toko_rating"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp;index:idx_toko_deleted_at"`

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
	GetActiveStores(limit, offset int, search string) ([]*Store, int64, error)
	GetPendingStores(limit, offset int, search string) ([]*Store, int64, error)
	GetActiveStoreByUserID(userID uint64) (*Store, error)
}

type CreateStoreRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	PhotoURL    string `json:"url_foto"`
}

type UpdateStoreRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string `json:"description,omitempty"`
	PhotoURL    *string `json:"url_foto,omitempty"`
}