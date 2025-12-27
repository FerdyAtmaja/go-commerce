package domain

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	Name       string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Detail     string         `json:"detail" gorm:"not null" validate:"required"`
	Phone      string         `json:"phone" gorm:"not null" validate:"required,min=10,max=15"`
	ProvinceID string         `json:"province_id" gorm:"not null" validate:"required"`
	CityID     string         `json:"city_id" gorm:"not null" validate:"required"`
	PostalCode string         `json:"postal_code" gorm:"not null" validate:"required"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type AddressRepository interface {
	Create(address *Address) error
	GetByID(id uint) (*Address, error)
	GetByUserID(userID uint, limit, offset int) ([]*Address, int64, error)
	Update(address *Address) error
	Delete(id uint) error
	CheckOwnership(addressID, userID uint) bool
}

type CreateAddressRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Detail     string `json:"detail" validate:"required"`
	Phone      string `json:"phone" validate:"required,min=10,max=15"`
	ProvinceID string `json:"province_id" validate:"required"`
	CityID     string `json:"city_id" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

type UpdateAddressRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Detail     string `json:"detail" validate:"required"`
	Phone      string `json:"phone" validate:"required,min=10,max=15"`
	ProvinceID string `json:"province_id" validate:"required"`
	CityID     string `json:"city_id" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required"`
}

// Province and City structures matching Indonesia API
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}