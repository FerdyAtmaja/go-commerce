package usecase

import (
	"errors"
	"fmt"
	"go-commerce/internal/domain"
	"time"
)

type TransactionUsecase struct {
	transactionRepo     domain.TransactionRepository
	transactionItemRepo domain.TransactionItemRepository
	productLogRepo      domain.ProductLogRepository
	productRepo         domain.ProductRepository
	addressRepo         domain.AddressRepository
	userRepo            domain.UserRepository
}

func NewTransactionUsecase(
	transactionRepo domain.TransactionRepository,
	transactionItemRepo domain.TransactionItemRepository,
	productLogRepo domain.ProductLogRepository,
	productRepo domain.ProductRepository,
	addressRepo domain.AddressRepository,
	userRepo domain.UserRepository,
) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo:     transactionRepo,
		transactionItemRepo: transactionItemRepo,
		productLogRepo:      productLogRepo,
		productRepo:         productRepo,
		addressRepo:         addressRepo,
		userRepo:            userRepo,
	}
}

func (u *TransactionUsecase) CreateTransaction(userID uint64, req *domain.CreateTransactionRequest) (*domain.Transaction, error) {
	// Validate address exists and belongs to user
	if !u.addressRepo.CheckOwnership(req.AlamatPengiriman, userID) {
		return nil, errors.New("address not found or access denied")
	}

	// Validate user exists
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Begin database transaction
	dbTx, err := u.transactionRepo.BeginTx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			u.transactionRepo.RollbackTx(dbTx)
			panic(r)
		}
	}()

	var totalAmount float64
	var transactionItems []*domain.TransactionItem

	// Validate products and calculate total
	for _, itemReq := range req.Items {
		product, err := u.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("product not found")
		}

		// Check stock availability
		if product.Stok < itemReq.Quantity {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("insufficient stock for product: " + product.NamaProduk)
		}

		// Create product log first
		productLog := &domain.ProductLog{
			ProductID:     itemReq.ProductID,
			NamaProduk:    product.NamaProduk,
			Slug:          product.Slug,
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Deskripsi:     product.Deskripsi,
			StoreID:       product.IDToko,
			CategoryID:    product.IDCategory,
		}

		// Create product log synchronously to get ID
		err = u.productLogRepo.Create(productLog)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, err
		}

		hargaSatuan := product.HargaKonsumen
		hargaTotal := hargaSatuan * float64(itemReq.Quantity)
		totalAmount += hargaTotal

		transactionItems = append(transactionItems, &domain.TransactionItem{
			ProductLogID:       productLog.ID,
			StoreID:            product.IDToko,
			Quantity:           itemReq.Quantity,
			HargaSatuan:        hargaSatuan,
			HargaTotal:         hargaTotal,
			NamaProdukSnapshot: product.NamaProduk,
		})
	}

	// Generate invoice code
	invoiceCode := fmt.Sprintf("INV-%d-%d", userID, time.Now().Unix())

	// Create transaction
	transaction := &domain.Transaction{
		UserID:           userID,
		AlamatPengiriman: req.AlamatPengiriman,
		HargaTotal:       totalAmount,
		KodeInvoice:      invoiceCode,
		MetodeBayar:      req.MetodeBayar,
		Status:           "pending",
	}

	err = u.transactionRepo.CreateWithTx(dbTx, transaction)
	if err != nil {
		u.transactionRepo.RollbackTx(dbTx)
		return nil, err
	}

	// Create transaction items
	for _, item := range transactionItems {
		item.TransactionID = transaction.ID
		
		err = u.transactionItemRepo.CreateWithTx(dbTx, item)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, err
		}
	}

	// Commit transaction
	err = u.transactionRepo.CommitTx(dbTx)
	if err != nil {
		u.transactionRepo.RollbackTx(dbTx)
		return nil, err
	}

	// Load relations for response
	transaction.User = user
	transaction.TransactionItems = transactionItems

	return transaction, nil
}

func (u *TransactionUsecase) GetTransactionByID(userID, transactionID uint64) (*domain.Transaction, error) {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if transaction.UserID != userID {
		return nil, errors.New("transaction not found or access denied")
	}

	return transaction, nil
}

func (u *TransactionUsecase) GetMyTransactions(userID uint64, page, limit int) ([]*domain.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.transactionRepo.GetByUserID(userID, limit, offset)
}