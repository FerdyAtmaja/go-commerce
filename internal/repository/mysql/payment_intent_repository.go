package mysql

import (
	"go-commerce/internal/domain"
	"gorm.io/gorm"
)

type paymentIntentRepository struct {
	db *gorm.DB
}

func NewPaymentIntentRepository(db *gorm.DB) domain.PaymentIntentRepository {
	return &paymentIntentRepository{db: db}
}

func (r *paymentIntentRepository) Create(intent *domain.PaymentIntent) error {
	return r.db.Create(intent).Error
}

func (r *paymentIntentRepository) GetByID(id uint) (*domain.PaymentIntent, error) {
	var intent domain.PaymentIntent
	err := r.db.Preload("Transaction").First(&intent, id).Error
	if err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *paymentIntentRepository) GetByTrxID(trxID uint) (*domain.PaymentIntent, error) {
	var intent domain.PaymentIntent
	err := r.db.Where("trx_id = ?", trxID).Preload("Transaction").First(&intent).Error
	if err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *paymentIntentRepository) UpdateStatus(id uint, status domain.PaymentIntentStatus) error {
	return r.db.Model(&domain.PaymentIntent{}).Where("id = ?", id).Update("status", status).Error
}

func (r *paymentIntentRepository) ExpireByTrxID(trxID uint) error {
	return r.db.Model(&domain.PaymentIntent{}).Where("trx_id = ? AND status = ?", trxID, domain.PaymentIntentStatusPending).Update("status", "expired").Error
}