package domain

import (
	"time"
)

type Transaction struct {
	ID                uint64    `json:"id" gorm:"primaryKey;autoIncrement;table:trx"`
	UserID            uint64    `json:"user_id" gorm:"column:id_user;not null"`
	AlamatPengiriman  uint64    `json:"alamat_pengiriman" gorm:"column:alamat_pengiriman;not null"`
	HargaTotal        float64   `json:"harga_total" gorm:"column:harga_total;not null"`
	KodeInvoice       string    `json:"kode_invoice" gorm:"column:kode_invoice;unique;not null"`
	MetodeBayar       string    `json:"metode_bayar" gorm:"column:metode_bayar"`
	Status            string    `json:"status" gorm:"column:status;not null;default:'pending'"`
	PaidAt            *time.Time `json:"paid_at" gorm:"column:paid_at"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	User            *User               `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Alamat          *Address            `json:"alamat,omitempty" gorm:"foreignKey:AlamatPengiriman;references:ID"`
	TransactionItems []*TransactionItem `json:"transaction_items,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
}

type TransactionItem struct {
	ID                   uint64    `json:"id" gorm:"primaryKey;autoIncrement;table:detail_trx"`
	TransactionID        uint64    `json:"transaction_id" gorm:"column:id_trx;not null"`
	ProductLogID         uint64    `json:"product_log_id" gorm:"column:id_log_produk;not null"`
	StoreID              uint64    `json:"store_id" gorm:"column:id_toko;not null"`
	Quantity             int       `json:"quantity" gorm:"column:kuantitas;not null"`
	HargaSatuan          float64   `json:"harga_satuan" gorm:"column:harga_satuan;not null"`
	HargaTotal           float64   `json:"harga_total" gorm:"column:harga_total;not null"`
	NamaProdukSnapshot   string    `json:"nama_produk_snapshot" gorm:"column:nama_produk_snapshot;not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
	ProductLog  *ProductLog  `json:"product_log,omitempty" gorm:"foreignKey:ProductLogID;references:ID"`
	Store       *Store       `json:"store,omitempty" gorm:"foreignKey:StoreID;references:ID"`
}

type ProductLog struct {
	ID             uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID      uint64    `json:"product_id" gorm:"column:id_produk;not null"`
	NamaProduk     string    `json:"nama_produk" gorm:"column:nama_produk;not null"`
	Slug           string    `json:"slug" gorm:"column:slug;not null"`
	HargaReseller  float64   `json:"harga_reseller" gorm:"column:harga_reseller;not null"`
	HargaKonsumen  float64   `json:"harga_konsumen" gorm:"column:harga_konsumen;not null"`
	Deskripsi      string    `json:"deskripsi" gorm:"column:deskripsi;type:text"`
	StoreID        uint64    `json:"store_id" gorm:"column:id_toko;not null"`
	CategoryID     uint64    `json:"category_id" gorm:"column:id_category;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Request DTOs
type CreateTransactionRequest struct {
	AlamatPengiriman uint64                      `json:"alamat_pengiriman" validate:"required"`
	MetodeBayar      string                      `json:"metode_bayar" validate:"required"`
	Items            []CreateTransactionItemRequest `json:"items" validate:"required,min=1"`
}

type CreateTransactionItemRequest struct {
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// Repository interfaces
type TransactionRepository interface {
	Create(tx *Transaction) error
	GetByID(id uint64) (*Transaction, error)
	GetByUserID(userID uint64, limit, offset int) ([]*Transaction, int64, error)
	Update(tx *Transaction) error
	BeginTx() (interface{}, error)
	CommitTx(tx interface{}) error
	RollbackTx(tx interface{}) error
	CreateWithTx(dbTx interface{}, tx *Transaction) error
}

type TransactionItemRepository interface {
	Create(item *TransactionItem) error
	CreateWithTx(dbTx interface{}, item *TransactionItem) error
	GetByTransactionID(transactionID uint64) ([]*TransactionItem, error)
}

type ProductLogRepository interface {
	Create(log *ProductLog) error
	CreateAsync(log *ProductLog)
	GetByProductID(productID uint64) ([]*ProductLog, error)
}