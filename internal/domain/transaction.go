package domain

import (
	"time"
)

type Transaction struct {
	ID                uint64     `json:"id" gorm:"primaryKey;column:id"`
	UserID            uint64     `json:"user_id" gorm:"column:id_user;type:bigint unsigned;not null;index:idx_trx_user"`
	AlamatPengiriman  uint64     `json:"alamat_pengiriman" gorm:"column:alamat_pengiriman;type:bigint unsigned;not null"`
	HargaTotal        float64    `json:"harga_total" gorm:"column:harga_total;type:decimal(14,2);not null"`
	KodeInvoice       string     `json:"kode_invoice" gorm:"column:kode_invoice;type:varchar(255);unique;not null;index:idx_trx_invoice"`
	MetodeBayar       string     `json:"metode_bayar" gorm:"column:metode_bayar;type:enum('transfer','cod','ewallet','credit_card')" validate:"omitempty,oneof=transfer cod ewallet credit_card"`
	Status            string     `json:"status_pembayaran" gorm:"column:status_pembayaran;type:enum('pending','paid','failed','refunded','cancelled','shipped','done');default:pending;index:idx_trx_status_pembayaran" validate:"omitempty,oneof=pending paid failed refunded cancelled shipped done"`
	OrderStatus       string     `json:"order_status" gorm:"column:order_status;type:enum('created','processed','shipped','delivered','cancelled');default:created;index:idx_trx_order_status" validate:"omitempty,oneof=created processed shipped delivered cancelled"`
	PaidAt            *time.Time `json:"paid_at" gorm:"column:paid_at;type:timestamp"`
	ShippedAt         *time.Time `json:"shipped_at" gorm:"column:shipped_at;type:timestamp"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_trx_created"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relations
	User            *User               `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Alamat          *Address            `json:"alamat,omitempty" gorm:"foreignKey:AlamatPengiriman;references:ID"`
	TransactionItems []*TransactionItem `json:"transaction_items,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
}

func (Transaction) TableName() string {
	return "trx"
}

type TransactionItem struct {
	ID                   uint64    `json:"id" gorm:"primaryKey;column:id"`
	TransactionID        uint64    `json:"transaction_id" gorm:"column:id_trx;type:bigint unsigned;not null;index:idx_detail_trx_trx"`
	ProductLogID         uint64    `json:"product_log_id" gorm:"column:id_log_produk;type:bigint unsigned;not null;index:idx_detail_trx_log_produk"`
	StoreID              uint64    `json:"store_id" gorm:"column:id_toko;type:bigint unsigned;not null;index:idx_detail_trx_toko"`
	Quantity             int       `json:"quantity" gorm:"column:kuantitas;type:int;not null"`
	HargaSatuan          float64   `json:"harga_satuan" gorm:"column:harga_satuan;type:decimal(12,2);not null"`
	HargaTotal           float64   `json:"harga_total" gorm:"column:harga_total;type:decimal(14,2);not null"`
	NamaProdukSnapshot   string    `json:"nama_produk_snapshot" gorm:"column:nama_produk_snapshot;type:varchar(255);not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relations
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
	ProductLog  *ProductLog  `json:"product_log,omitempty" gorm:"foreignKey:ProductLogID;references:ID"`
	Store       *Store       `json:"store,omitempty" gorm:"foreignKey:StoreID;references:ID"`
}

func (TransactionItem) TableName() string {
	return "detail_trx"
}

type ProductLog struct {
	ID             uint64    `json:"id" gorm:"primaryKey;column:id"`
	ProductID      uint64    `json:"product_id" gorm:"column:id_produk;type:bigint unsigned;not null;index:idx_log_produk_produk"`
	NamaProduk     string    `json:"nama_produk" gorm:"column:nama_produk;type:varchar(255);not null"`
	Slug           string    `json:"slug" gorm:"column:slug;type:varchar(255);not null"`
	HargaReseller  float64   `json:"harga_reseller" gorm:"column:harga_reseller;type:decimal(12,2);not null"`
	HargaKonsumen  float64   `json:"harga_konsumen" gorm:"column:harga_konsumen;type:decimal(12,2);not null"`
	Deskripsi      string    `json:"deskripsi" gorm:"column:deskripsi;type:text"`
	StoreID        uint64    `json:"store_id" gorm:"column:id_toko;type:bigint;not null;index:idx_log_produk_toko"`
	CategoryID     uint64    `json:"category_id" gorm:"column:id_category;type:bigint;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_log_produk_created"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (ProductLog) TableName() string {
	return "log_produk"
}

// Request DTOs
type CreateTransactionRequest struct {
	AlamatPengiriman uint64                        `json:"alamat_pengiriman" validate:"required"`
	MetodeBayar      string                        `json:"metode_bayar" validate:"required,oneof=transfer cod ewallet credit_card"`
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
	GetByStatus(status string, limit, offset int) ([]*Transaction, int64, error)
	Update(tx *Transaction) error
	UpdateStatus(id uint64, status string) error
	UpdateStatusWithTx(dbTx interface{}, id uint64, status string) error
	UpdateOrderStatus(id uint64, orderStatus string) error
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
	GetByID(id uint64) (*ProductLog, error)
	GetByProductID(productID uint64) ([]*ProductLog, error)
}