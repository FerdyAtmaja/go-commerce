package domain

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID            uint64         `json:"id" gorm:"primaryKey"`
	UserID        uint64         `json:"user_id" gorm:"column:id_user;not null;index"`
	JudulAlamat   string         `json:"judul_alamat" gorm:"column:judul_alamat;not null" validate:"required,min=2,max=255"`
	NamaPenerima  string         `json:"nama_penerima" gorm:"column:nama_penerima;not null" validate:"required,min=2,max=255"`
	NoTelp        string         `json:"notelp" gorm:"column:notelp" validate:"omitempty,min=10,max=20"`
	DetailAlamat  string         `json:"detail_alamat" gorm:"column:detail_alamat;type:text;not null" validate:"required"`
	KodePos       string         `json:"kode_pos" gorm:"column:kode_pos" validate:"omitempty,max=10"`
	IsDefault     bool           `json:"is_default" gorm:"column:is_default;default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (Address) TableName() string {
	return "alamat"
}

type AddressRepository interface {
	Create(address *Address) error
	GetByID(id uint64) (*Address, error)
	GetByUserID(userID uint64, limit, offset int) ([]*Address, int64, error)
	GetDefaultByUserID(userID uint64) (*Address, error)
	Update(address *Address) error
	Delete(id uint64) error
	SetDefault(addressID, userID uint64) error
	CheckOwnership(addressID, userID uint64) bool
}

type CreateAddressRequest struct {
	JudulAlamat  string `json:"judul_alamat" validate:"required,min=2,max=255"`
	NamaPenerima string `json:"nama_penerima" validate:"required,min=2,max=255"`
	NoTelp       string `json:"notelp" validate:"omitempty,min=10,max=20"`
	DetailAlamat string `json:"detail_alamat" validate:"required"`
	KodePos      string `json:"kode_pos" validate:"omitempty,max=10"`
	IsDefault    bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	JudulAlamat  string `json:"judul_alamat" validate:"required,min=2,max=255"`
	NamaPenerima string `json:"nama_penerima" validate:"required,min=2,max=255"`
	NoTelp       string `json:"notelp" validate:"omitempty,min=10,max=20"`
	DetailAlamat string `json:"detail_alamat" validate:"required"`
	KodePos      string `json:"kode_pos" validate:"omitempty,max=10"`
	IsDefault    bool   `json:"is_default"`
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