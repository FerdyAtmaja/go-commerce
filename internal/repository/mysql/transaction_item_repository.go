package mysql

import (
	"go-commerce/internal/domain"

	"gorm.io/gorm"
)

type transactionItemRepository struct {
	db *gorm.DB
}

func NewTransactionItemRepository(db *gorm.DB) domain.TransactionItemRepository {
	return &transactionItemRepository{db: db}
}

func (r *transactionItemRepository) Create(item *domain.TransactionItem) error {
	return r.db.Create(item).Error
}

func (r *transactionItemRepository) CreateWithTx(dbTx interface{}, item *domain.TransactionItem) error {
	gormTx := dbTx.(*gorm.DB)
	return gormTx.Create(item).Error
}

func (r *transactionItemRepository) GetByTransactionID(transactionID uint64) ([]*domain.TransactionItem, error) {
	var items []*domain.TransactionItem
	err := r.db.Where("id_trx = ?", transactionID).Find(&items).Error
	return items, err
}