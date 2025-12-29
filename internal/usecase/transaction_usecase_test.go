package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionUsecase_GetTransactionByID_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockTransactionItemRepo := new(mocks.MockTransactionItemRepository)
	mockProductLogRepo := new(mocks.MockProductLogRepository)
	mockProductRepo := new(mocks.ProductRepositoryMock)
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.StoreRepositoryMock)

	transactionUsecase := NewTransactionUsecase(
		mockTransactionRepo,
		mockTransactionItemRepo,
		mockProductLogRepo,
		mockProductRepo,
		mockAddressRepo,
		mockUserRepo,
		mockStoreRepo,
	)

	userID := uint64(1)
	transactionID := uint64(1)

	transaction := &domain.Transaction{
		ID:     transactionID,
		UserID: userID,
		Status: "pending",
	}

	// Mock expectations
	mockTransactionRepo.On("GetByID", transactionID).Return(transaction, nil)

	// Execute
	result, err := transactionUsecase.GetTransactionByID(userID, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, userID, result.UserID)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionUsecase_GetTransactionByID_AccessDenied(t *testing.T) {
	// Setup
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockTransactionItemRepo := new(mocks.MockTransactionItemRepository)
	mockProductLogRepo := new(mocks.MockProductLogRepository)
	mockProductRepo := new(mocks.ProductRepositoryMock)
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.StoreRepositoryMock)

	transactionUsecase := NewTransactionUsecase(
		mockTransactionRepo,
		mockTransactionItemRepo,
		mockProductLogRepo,
		mockProductRepo,
		mockAddressRepo,
		mockUserRepo,
		mockStoreRepo,
	)

	userID := uint64(1)
	transactionID := uint64(1)

	transaction := &domain.Transaction{
		ID:     transactionID,
		UserID: 2, // Different user
		Status: "pending",
	}

	// Mock expectations
	mockTransactionRepo.On("GetByID", transactionID).Return(transaction, nil)

	// Execute
	result, err := transactionUsecase.GetTransactionByID(userID, transactionID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "access denied")

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionUsecase_GetMyTransactions_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockTransactionItemRepo := new(mocks.MockTransactionItemRepository)
	mockProductLogRepo := new(mocks.MockProductLogRepository)
	mockProductRepo := new(mocks.ProductRepositoryMock)
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.StoreRepositoryMock)

	transactionUsecase := NewTransactionUsecase(
		mockTransactionRepo,
		mockTransactionItemRepo,
		mockProductLogRepo,
		mockProductRepo,
		mockAddressRepo,
		mockUserRepo,
		mockStoreRepo,
	)

	userID := uint64(1)
	page := 1
	limit := 10
	offset := 0

	transactions := []*domain.Transaction{
		{ID: 1, UserID: userID, Status: "pending"},
		{ID: 2, UserID: userID, Status: "paid"},
	}
	total := int64(2)

	// Mock expectations
	mockTransactionRepo.On("GetByUserID", userID, limit, offset).Return(transactions, total, nil)

	// Execute
	result, resultTotal, err := transactionUsecase.GetMyTransactions(userID, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, total, resultTotal)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionUsecase_CreateTransaction_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockTransactionItemRepo := new(mocks.MockTransactionItemRepository)
	mockProductLogRepo := new(mocks.MockProductLogRepository)
	mockProductRepo := new(mocks.ProductRepositoryMock)
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.StoreRepositoryMock)

	transactionUsecase := NewTransactionUsecase(
		mockTransactionRepo,
		mockTransactionItemRepo,
		mockProductLogRepo,
		mockProductRepo,
		mockAddressRepo,
		mockUserRepo,
		mockStoreRepo,
	)

	userID := uint64(1)
	addressID := uint64(1)
	productID := uint64(1)

	req := &domain.CreateTransactionRequest{
		AlamatPengiriman: addressID,
		MetodeBayar:      "transfer",
		Items: []domain.CreateTransactionItemRequest{
			{ProductID: productID, Quantity: 2},
		},
	}

	user := &domain.User{ID: userID, Name: "Test User"}
	product := &domain.Product{
		ID:             productID,
		NamaProduk:     "Test Product",
		Slug:           "test-product",
		HargaKonsumen:  10000.0,
		Stok:           10,
		IDToko:         1,
		IDCategory:     1,
		Status:         "active",
	}

	store := &domain.Store{
		ID:     1,
		UserID: 2,
		Status: "active",
	}

	mockTx := "mock_transaction"

	// Mock expectations - Atomic operations
	mockAddressRepo.On("CheckOwnership", addressID, userID).Return(true)
	mockUserRepo.On("GetByID", userID).Return(user, nil)
	mockTransactionRepo.On("BeginTx").Return(mockTx, nil)
	mockTransactionRepo.On("RollbackTx", mockTx).Return(nil).Maybe() // Add this for defer
	mockProductRepo.On("GetByID", productID).Return(product, nil)
	mockStoreRepo.On("GetByID", uint64(1)).Return(store, nil)
	mockProductLogRepo.On("Create", mock.AnythingOfType("*domain.ProductLog")).Return(nil)
	mockTransactionRepo.On("CreateWithTx", mockTx, mock.AnythingOfType("*domain.Transaction")).Return(nil)
	mockTransactionItemRepo.On("CreateWithTx", mockTx, mock.AnythingOfType("*domain.TransactionItem")).Return(nil)
	mockTransactionRepo.On("CommitTx", mockTx).Return(nil)

	// Execute
	result, err := transactionUsecase.CreateTransaction(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, addressID, result.AlamatPengiriman)
	assert.Equal(t, "transfer", result.MetodeBayar)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, "created", result.OrderStatus)
	assert.Equal(t, 20000.0, result.HargaTotal) // 2 * 10000

	// Verify all mocks called
	mockAddressRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockProductLogRepo.AssertExpectations(t)
	mockTransactionItemRepo.AssertExpectations(t)
}