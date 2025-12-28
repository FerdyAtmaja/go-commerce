package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID              uint64         `json:"id" gorm:"primaryKey"`
	NamaProduk      string         `json:"nama_produk" gorm:"not null" validate:"required,min=2,max=255"`
	Slug            string         `json:"slug" gorm:"uniqueIndex;not null"`
	HargaReseller   float64        `json:"harga_reseller" gorm:"not null" validate:"required,min=0"`
	HargaKonsumen   float64        `json:"harga_konsumen" gorm:"not null" validate:"required,min=0"`
	Stok            int            `json:"stok" gorm:"default:0"`
	Deskripsi       string         `json:"deskripsi"`
	IDToko          uint64         `json:"id_toko" gorm:"not null;index"`
	IDCategory      uint64         `json:"id_category" gorm:"not null;index"`
	Status          string         `json:"status" gorm:"default:active"`
	Berat           int            `json:"berat" gorm:"default:0"`
	SoldCount       int            `json:"sold_count" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Toko     Store        `json:"toko,omitempty" gorm:"foreignKey:IDToko"`
	Category Category     `json:"category,omitempty" gorm:"foreignKey:IDCategory"`
	Photos   []PhotoProduk `json:"photos,omitempty" gorm:"foreignKey:IDProduk"`
}

func (Product) TableName() string {
	return "produk"
}

type PhotoProduk struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	IDProduk  uint64         `json:"id_produk" gorm:"not null;index"`
	URL       string         `json:"url" gorm:"not null"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	Position  int            `json:"position" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Product Product `json:"product,omitempty" gorm:"foreignKey:IDProduk"`
}

func (PhotoProduk) TableName() string {
	return "foto_produk"
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id uint64) (*Product, error)
	GetBySlug(slug string) (*Product, error)
	GetByTokoID(tokoID uint64, limit, offset int, search string) ([]*Product, int64, error)
	GetAll(limit, offset int, search, categoryID string) ([]*Product, int64, error)
	Update(product *Product) error
	Delete(id uint64) error
	CheckOwnership(productID, tokoID uint64) error
}

type PhotoProdukRepository interface {
	Create(photo *PhotoProduk) error
	GetByProductID(productID uint64) ([]*PhotoProduk, error)
	Update(photo *PhotoProduk) error
	Delete(id uint64) error
	SetPrimary(productID, photoID uint64) error
}

type CreateProductRequest struct {
	NamaProduk    string  `json:"nama_produk" validate:"required,min=2,max=255"`
	HargaReseller float64 `json:"harga_reseller" validate:"required,min=0"`
	HargaKonsumen float64 `json:"harga_konsumen" validate:"required,min=0"`
	Stok          int     `json:"stok" validate:"min=0"`
	Deskripsi     string  `json:"deskripsi"`
	IDCategory    uint64  `json:"id_category" validate:"required"`
	Berat         int     `json:"berat" validate:"min=0"`
}

type UpdateProductRequest struct {
	NamaProduk    string  `json:"nama_produk" validate:"required,min=2,max=255"`
	HargaReseller float64 `json:"harga_reseller" validate:"required,min=0"`
	HargaKonsumen float64 `json:"harga_konsumen" validate:"required,min=0"`
	Stok          int     `json:"stok" validate:"min=0"`
	Deskripsi     string  `json:"deskripsi"`
	IDCategory    uint64  `json:"id_category" validate:"required"`
	Berat         int     `json:"berat" validate:"min=0"`
}

type ProductFilter struct {
	Search     string `json:"search"`
	CategoryID string `json:"category_id"`
	MinPrice   string `json:"min_price"`
	MaxPrice   string `json:"max_price"`
	SortBy     string `json:"sort_by"` // price_asc, price_desc, newest, oldest, popular
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}