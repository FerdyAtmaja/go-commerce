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
	storeRepo           domain.StoreRepository
}

func NewTransactionUsecase(
	transactionRepo domain.TransactionRepository,
	transactionItemRepo domain.TransactionItemRepository,
	productLogRepo domain.ProductLogRepository,
	productRepo domain.ProductRepository,
	addressRepo domain.AddressRepository,
	userRepo domain.UserRepository,
	storeRepo domain.StoreRepository,
) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo:     transactionRepo,
		transactionItemRepo: transactionItemRepo,
		productLogRepo:      productLogRepo,
		productRepo:         productRepo,
		addressRepo:         addressRepo,
		userRepo:            userRepo,
		storeRepo:           storeRepo,
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
			return nil, errors.New("product not found or not available")
		}

		// Get store and validate status
		store, err := u.storeRepo.GetByID(product.IDToko)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("store not found")
		}

		// Check if store is active
		if store.Status != "active" {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("STORE_NOT_AVAILABLE")
		}

		// Check if product is active
		if product.Status != "active" {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("PRODUCT_NOT_AVAILABLE")
		}

		// Check stock availability
		if product.Stok < itemReq.Quantity {
			u.transactionRepo.RollbackTx(dbTx)
			return nil, errors.New("INSUFFICIENT_STOCK")
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
		OrderStatus:      "created",
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

func (u *TransactionUsecase) UpdatePaymentStatus(userID, transactionID uint64, status string) error {
	// Validate status
	validStatuses := []string{"pending", "paid", "cancelled", "shipped", "done"}
	validStatus := false
	for _, s := range validStatuses {
		if s == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return errors.New("invalid status")
	}

	// Check ownership
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}
	if transaction.UserID != userID {
		return errors.New("access denied")
	}

	return u.transactionRepo.UpdateStatus(transactionID, status)
}

// OnPaymentPaid - Used by payment intent system
func (u *TransactionUsecase) OnPaymentPaid(transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return err
	}

	// Check if transaction is cancelled
	if transaction.Status == "cancelled" || transaction.OrderStatus == "cancelled" {
		return errors.New("transaction cancelled")
	}

	// Idempotent check
	if transaction.Status != "pending" {
		return nil // Already processed
	}

	// Begin database transaction
	dbTx, err := u.transactionRepo.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			u.transactionRepo.RollbackTx(dbTx)
			panic(r)
		}
	}()

	// Update payment status
	err = u.transactionRepo.UpdateStatus(transactionID, "paid")
	if err != nil {
		u.transactionRepo.RollbackTx(dbTx)
		return err
	}

	// Reduce stock and update sold count for each item
	for _, item := range transaction.TransactionItems {
		// Get product ID from product log
		productLogs, err := u.productLogRepo.GetByProductID(item.ProductLogID)
		if err != nil || len(productLogs) == 0 {
			u.transactionRepo.RollbackTx(dbTx)
			return errors.New("product log not found")
		}

		productID := productLogs[0].ProductID

		// Update stock
		err = u.productRepo.UpdateStockWithTx(dbTx, productID, item.Quantity)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return err
		}

		// Update sold count
		err = u.productRepo.UpdateSoldCountWithTx(dbTx, productID, item.Quantity)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return err
		}
	}

	return u.transactionRepo.CommitTx(dbTx)
}

// OnPaymentFailed - Used by payment intent system
func (u *TransactionUsecase) OnPaymentFailed(transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return err
	}

	// Idempotent check
	if transaction.Status != "pending" {
		return nil
	}

	// Update status
	err = u.transactionRepo.UpdateStatus(transactionID, "failed")
	if err != nil {
		return err
	}

	return u.transactionRepo.UpdateOrderStatus(transactionID, "cancelled")
}

// ProcessOrder - Seller processes order
func (u *TransactionUsecase) ProcessOrder(sellerID, transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Check if transaction is cancelled
	if transaction.Status == "cancelled" || transaction.OrderStatus == "cancelled" {
		return errors.New("transaction cancelled")
	}

	// Validate seller owns store in transaction
	for _, item := range transaction.TransactionItems {
		store, err := u.storeRepo.GetByID(item.StoreID)
		if err != nil || store.UserID != sellerID {
			return errors.New("forbidden: seller does not own store")
		}
	}

	// Validate payment status
	if transaction.Status != "paid" {
		return errors.New("payment not completed")
	}

	// Validate order status
	if transaction.OrderStatus != "created" {
		return errors.New("invalid state")
	}

	return u.transactionRepo.UpdateOrderStatus(transactionID, "processed")
}

// ShipOrder - Seller ships order
func (u *TransactionUsecase) ShipOrder(sellerID, transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Validate seller owns store
	for _, item := range transaction.TransactionItems {
		store, err := u.storeRepo.GetByID(item.StoreID)
		if err != nil || store.UserID != sellerID {
			return errors.New("forbidden")
		}
	}

	// Validate order status
	if transaction.OrderStatus != "processed" {
		return errors.New("order not ready")
	}

	return u.transactionRepo.UpdateOrderStatus(transactionID, "shipped")
}

// ConfirmDelivered - Buyer confirms delivery
func (u *TransactionUsecase) ConfirmDelivered(userID, transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Check ownership
	if transaction.UserID != userID {
		return errors.New("forbidden")
	}

	// Validate order status
	if transaction.OrderStatus != "shipped" {
		return errors.New("order not shipped")
	}

	return u.transactionRepo.UpdateOrderStatus(transactionID, "delivered")
}

// CancelTransaction - Buyer cancels transaction
func (u *TransactionUsecase) CancelTransaction(userID, transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Check ownership
	if transaction.UserID != userID {
		return errors.New("forbidden")
	}

	// Check if paid - should use refund instead
	if transaction.Status == "paid" {
		return errors.New("use refund for paid transactions")
	}

	// Can only cancel created orders (not processed, shipped, or delivered)
	if transaction.OrderStatus != "created" {
		return errors.New("cannot cancel processed order")
	}

	// Update both status and order status to cancelled
	err = u.transactionRepo.UpdateStatus(transactionID, "cancelled")
	if err != nil {
		return err
	}

	err = u.transactionRepo.UpdateOrderStatus(transactionID, "cancelled")
	if err != nil {
		return err
	}

	return nil
}

// RefundTransaction - Admin/System refunds transaction
func (u *TransactionUsecase) RefundTransaction(transactionID uint64) error {
	transaction, err := u.transactionRepo.GetByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Can only refund paid transactions
	if transaction.Status != "paid" {
		return errors.New("NOT_REFUNDABLE")
	}

	// Begin database transaction
	dbTx, err := u.transactionRepo.BeginTx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			u.transactionRepo.RollbackTx(dbTx)
			panic(r)
		}
	}()

	// Update payment status to refunded
	err = u.transactionRepo.UpdateStatus(transactionID, "refunded")
	if err != nil {
		u.transactionRepo.RollbackTx(dbTx)
		return err
	}

	// Update order status to cancelled
	err = u.transactionRepo.UpdateOrderStatus(transactionID, "cancelled")
	if err != nil {
		u.transactionRepo.RollbackTx(dbTx)
		return err
	}

	// Restore stock and reduce sold count for each item
	for _, item := range transaction.TransactionItems {
		// Get product ID from product log
		productLogs, err := u.productLogRepo.GetByProductID(item.ProductLogID)
		if err != nil || len(productLogs) == 0 {
			u.transactionRepo.RollbackTx(dbTx)
			return errors.New("product log not found")
		}

		productID := productLogs[0].ProductID

		// Restore stock (add back)
		err = u.productRepo.UpdateStockWithTx(dbTx, productID, -item.Quantity)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return err
		}

		// Reduce sold count
		err = u.productRepo.UpdateSoldCountWithTx(dbTx, productID, -item.Quantity)
		if err != nil {
			u.transactionRepo.RollbackTx(dbTx)
			return err
		}
	}

	return u.transactionRepo.CommitTx(dbTx)
}