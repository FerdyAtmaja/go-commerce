package mysql

import (
	"go-commerce/internal/domain"
	"time"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetByID(id uint64) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.Preload("User").Preload("Alamat").Preload("TransactionItems").First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByUserID(userID uint64, limit, offset int) ([]*domain.Transaction, int64, error) {
	var transactions []*domain.Transaction
	var total int64

	// Use idx_trx_user index
	err := r.db.Model(&domain.Transaction{}).Where("id_user = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("id_user = ?", userID).
		Preload("Alamat").
		Preload("TransactionItems").
		Order("created_at DESC"). // Uses idx_trx_created index
		Limit(limit).Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) Update(tx *domain.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *transactionRepository) UpdateStatus(id uint64, status string) error {
	updates := map[string]interface{}{"status_pembayaran": status}
	if status == "paid" {
		updates["paid_at"] = time.Now()
	}
	return r.db.Model(&domain.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *transactionRepository) UpdatePaymentStatus(id uint64, paymentStatus string) error {
	updates := map[string]interface{}{"payment_status": paymentStatus}
	if paymentStatus == "paid" {
		updates["paid_at"] = time.Now()
	}
	return r.db.Model(&domain.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *transactionRepository) UpdateOrderStatus(id uint64, orderStatus string) error {
	updates := map[string]interface{}{"order_status": orderStatus}
	if orderStatus == "shipped" {
		updates["shipped_at"] = time.Now()
	}
	return r.db.Model(&domain.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *transactionRepository) BeginTx() (interface{}, error) {
	return r.db.Begin(), nil
}

func (r *transactionRepository) CommitTx(tx interface{}) error {
	dbTx := tx.(*gorm.DB)
	return dbTx.Commit().Error
}

func (r *transactionRepository) RollbackTx(tx interface{}) error {
	dbTx := tx.(*gorm.DB)
	return dbTx.Rollback().Error
}

func (r *transactionRepository) CreateWithTx(dbTx interface{}, tx *domain.Transaction) error {
	gormTx := dbTx.(*gorm.DB)
	return gormTx.Create(tx).Error
}

// GetByStatus gets transactions by status (uses idx_trx_status_pembayaran index)
func (r *transactionRepository) GetByStatus(status string, limit, offset int) ([]*domain.Transaction, int64, error) {
	var transactions []*domain.Transaction
	var total int64

	err := r.db.Model(&domain.Transaction{}).Where("status_pembayaran = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status_pembayaran = ?", status).
		Preload("User").
		Preload("Alamat").
		Preload("TransactionItems").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}