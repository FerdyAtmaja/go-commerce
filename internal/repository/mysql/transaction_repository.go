package mysql

import (
	"go-commerce/internal/domain"

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
	err := r.db.Preload("User").Preload("Store").Preload("TransactionItems").First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByUserID(userID uint64, limit, offset int) ([]*domain.Transaction, int64, error) {
	var transactions []*domain.Transaction
	var total int64

	err := r.db.Model(&domain.Transaction{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).
		Preload("Store").
		Preload("TransactionItems").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) GetByStoreID(storeID uint64, limit, offset int) ([]*domain.Transaction, int64, error) {
	var transactions []*domain.Transaction
	var total int64

	err := r.db.Model(&domain.Transaction{}).Where("id_toko = ?", storeID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("id_toko = ?", storeID).
		Preload("User").
		Preload("TransactionItems").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) Update(tx *domain.Transaction) error {
	return r.db.Save(tx).Error
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