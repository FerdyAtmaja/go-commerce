package domain

import (
	"time"
)

type PaymentIntentStatus string

const (
	PaymentIntentStatusPending PaymentIntentStatus = "pending"
	PaymentIntentStatusSuccess PaymentIntentStatus = "success"
	PaymentIntentStatusFailed  PaymentIntentStatus = "failed"
)

type PaymentIntent struct {
	ID        uint                `json:"id" gorm:"primaryKey"`
	TrxID     uint                `json:"trx_id" gorm:"not null;index"`
	Method    string              `json:"method" gorm:"size:50;not null"`
	Status    PaymentIntentStatus `json:"status" gorm:"size:20;not null;default:'pending'"`
	ExpiredAt time.Time           `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	
	// Relations
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TrxID"`
}

type PaymentIntentRepository interface {
	Create(intent *PaymentIntent) error
	GetByID(id uint) (*PaymentIntent, error)
	GetByTrxID(trxID uint) (*PaymentIntent, error)
	UpdateStatus(id uint, status PaymentIntentStatus) error
	ExpireByTrxID(trxID uint) error
}

type PaymentIntentUsecase interface {
	CreatePaymentIntent(trxID uint, method string) (*PaymentIntent, error)
	ProcessPaymentSuccess(intentID uint) error
	ProcessPaymentFailed(intentID uint) error
	ExpireIntentsByTrxID(trxID uint) error
}