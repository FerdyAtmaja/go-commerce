package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID              uint64         `json:"id" gorm:"primaryKey;column:id"`
	NamaProduk      string         `json:"nama_produk" gorm:"column:nama_produk;type:varchar(255);not null" validate:"required,min=2,max=255"`
	Slug            string         `json:"slug" gorm:"column:slug;type:varchar(255);uniqueIndex;not null"`
	HargaReseller   float64        `json:"harga_reseller" gorm:"column:harga_reseller;type:decimal(12,2);not null" validate:"required,min=0"`
	HargaKonsumen   float64        `json:"harga_konsumen" gorm:"column:harga_konsumen;type:decimal(12,2);not null" validate:"required,min=0"`
	Stok            int            `json:"stok" gorm:"column:stok;type:int;default:0"`
	Deskripsi       string         `json:"deskripsi" gorm:"column:deskripsi;type:text"`
	IDToko          uint64         `json:"id_toko" gorm:"column:id_toko;type:bigint unsigned;not null;index:idx_produk_toko"`
	IDCategory      uint64         `json:"id_category" gorm:"column:id_category;type:bigint unsigned;not null;index:idx_produk_category"`
	Status          string         `json:"status" gorm:"column:status;type:enum('active','inactive');default:active;index:idx_produk_status" validate:"oneof=active inactive"`
	Berat           int            `json:"berat" gorm:"column:berat;type:int;default:0"`
	SoldCount       int            `json:"sold_count" gorm:"column:sold_count;type:int;default:0"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp;index:idx_produk_deleted_at"`

	// Relations
	Toko     Store        `json:"toko,omitempty" gorm:"foreignKey:IDToko"`
	Category Category     `json:"category,omitempty" gorm:"foreignKey:IDCategory"`
	Photos   []PhotoProduk `json:"photos,omitempty" gorm:"foreignKey:IDProduk"`
}

func (Product) TableName() string {
	return "produk"
}

type PhotoProduk struct {
	ID        uint64         `json:"id" gorm:"primaryKey;column:id"`
	IDProduk  uint64         `json:"id_produk" gorm:"column:id_produk;type:bigint unsigned;not null;index:idx_foto_produk_produk"`
	URL       string         `json:"url" gorm:"column:url;type:varchar(255);not null"`
	IsPrimary bool           `json:"is_primary" gorm:"column:is_primary;type:boolean;default:false;index:idx_foto_produk_primary"`
	Position  int64          `json:"position" gorm:"column:position;type:bigint;default:0;index:idx_foto_produk_position"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp;index:idx_foto_produk_deleted_at"`

	// Relations
	Product Product `json:"product,omitempty" gorm:"foreignKey:IDProduk"`
}

func (PhotoProduk) TableName() string {
	return "foto_produk"
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id uint64) (*Product, error)
	GetByIDForManagement(id uint64) (*Product, error)
	GetBySlug(slug string) (*Product, error)
	SearchBySlug(slugPattern string, limit, offset int) ([]*Product, int64, error)
	GetByTokoID(tokoID uint64, limit, offset int, search string) ([]*Product, int64, error)
	GetAll(limit, offset int, search, categoryID string) ([]*Product, int64, error)
	GetAllWithFilter(filter *ProductFilter) ([]*Product, int64, error)
	GetByStatus(status string, limit, offset int) ([]*Product, int64, error)
	Update(product *Product) error
	GetStockWithLock(dbTx interface{}, productID uint64) (int, error)
	UpdateStockWithTx(dbTx interface{}, productID uint64, quantity int) error
	UpdateSoldCountWithTx(dbTx interface{}, productID uint64, quantity int) error
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
	Status        string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

type UpdateProductRequest struct {
	NamaProduk    *string  `json:"nama_produk,omitempty" validate:"omitempty,min=2,max=255"`
	HargaReseller *float64 `json:"harga_reseller,omitempty" validate:"omitempty,min=0"`
	HargaKonsumen *float64 `json:"harga_konsumen,omitempty" validate:"omitempty,min=0"`
	Stok          *int     `json:"stok,omitempty" validate:"omitempty,min=0"`
	Deskripsi     *string  `json:"deskripsi,omitempty"`
	IDCategory    *uint64  `json:"id_category,omitempty" validate:"omitempty"`
	Berat         *int     `json:"berat,omitempty" validate:"omitempty,min=0"`
	Status        *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
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