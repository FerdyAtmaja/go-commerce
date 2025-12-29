package usecase

import (
	"errors"
	"go-commerce/internal/domain"
	"time"
)

type paymentIntentUsecase struct {
	paymentIntentRepo domain.PaymentIntentRepository
	transactionRepo   domain.TransactionRepository
	transactionUC     *TransactionUsecase
}

func NewPaymentIntentUsecase(
	paymentIntentRepo domain.PaymentIntentRepository,
	transactionRepo domain.TransactionRepository,
	transactionUC *TransactionUsecase,
) domain.PaymentIntentUsecase {
	return &paymentIntentUsecase{
		paymentIntentRepo: paymentIntentRepo,
		transactionRepo:   transactionRepo,
		transactionUC:     transactionUC,
	}
}

func (uc *paymentIntentUsecase) CreatePaymentIntent(trxID uint, method string) (*domain.PaymentIntent, error) {
	// Check if transaction exists and is pending
	trx, err := uc.transactionRepo.GetByID(uint64(trxID))
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	
	if trx.PaymentStatus != "pending" {
		return nil, errors.New("transaction is not in pending status")
	}
	
	// Check if payment intent already exists
	existing, _ := uc.paymentIntentRepo.GetByTrxID(trxID)
	if existing != nil {
		return existing, nil // Return existing intent (idempotent)
	}
	
	// Create new payment intent with 30 minutes expiry
	intent := &domain.PaymentIntent{
		TrxID:     trxID,
		Method:    method,
		Status:    domain.PaymentIntentStatusPending,
		ExpiredAt: time.Now().Add(30 * time.Minute),
	}
	
	err = uc.paymentIntentRepo.Create(intent)
	if err != nil {
		return nil, err
	}
	
	return intent, nil
}

func (uc *paymentIntentUsecase) ProcessPaymentSuccess(intentID uint) error {
	intent, err := uc.paymentIntentRepo.GetByID(intentID)
	if err != nil {
		return errors.New("payment intent not found")
	}
	
	// Decision table logic
	switch {
	case intent.Status == domain.PaymentIntentStatusSuccess && intent.Transaction.PaymentStatus == "paid":
		// Idempotent - already processed
		return nil
		
	case intent.Status == domain.PaymentIntentStatusPending && intent.Transaction.PaymentStatus == "pending":
		// Process payment success
		err = uc.paymentIntentRepo.UpdateStatus(intentID, domain.PaymentIntentStatusSuccess)
		if err != nil {
			return err
		}
		
		// Update transaction to paid
		return uc.transactionUC.OnPaymentPaid(uint64(intent.TrxID))
		
	default:
		return errors.New("invalid payment intent state for success processing")
	}
}

func (uc *paymentIntentUsecase) ProcessPaymentFailed(intentID uint) error {
	intent, err := uc.paymentIntentRepo.GetByID(intentID)
	if err != nil {
		return errors.New("payment intent not found")
	}
	
	if intent.Status != domain.PaymentIntentStatusPending {
		return nil // Idempotent
	}
	
	// Update intent status
	err = uc.paymentIntentRepo.UpdateStatus(intentID, domain.PaymentIntentStatusFailed)
	if err != nil {
		return err
	}
	
	// Update transaction to failed
	return uc.transactionUC.OnPaymentFailed(uint64(intent.TrxID))
}

func (uc *paymentIntentUsecase) ExpireIntentsByTrxID(trxID uint) error {
	return uc.paymentIntentRepo.ExpireByTrxID(trxID)
}