package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestStoreUsecase_GetMyStore_Success(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	userID := uint64(1)
	store := &domain.Store{
		ID:          1,
		UserID:      userID,
		Name:        "Test Store",
		Description: "Test Description",
	}

	// Mock expectations
	mockStoreRepo.On("GetByUserID", userID).Return(store, nil)

	// Execute
	result, err := storeUsecase.GetMyStore(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, store.ID, result.ID)
	assert.Equal(t, store.UserID, result.UserID)
	assert.Equal(t, store.Name, result.Name)

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_GetMyStore_NotFound(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	userID := uint64(999)

	// Mock expectations
	mockStoreRepo.On("GetByUserID", userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := storeUsecase.GetMyStore(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "STORE_NOT_FOUND")

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_UpdateMyStore_Success(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	userID := uint64(1)
	existingStore := &domain.Store{
		ID:          1,
		UserID:      userID,
		Name:        "Old Store Name",
		Description: "Old Description",
	}

	newName := "New Store Name"
	newDescription := "New Description"
	req := &domain.UpdateStoreRequest{
		Name:        &newName,
		Description: &newDescription,
	}

	// Mock expectations
	mockStoreRepo.On("GetByUserID", userID).Return(existingStore, nil)
	mockStoreRepo.On("Update", existingStore).Return(nil)

	// Execute
	result, err := storeUsecase.UpdateMyStore(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, *req.Name, result.Name)
	assert.Equal(t, *req.Description, result.Description)
	assert.Equal(t, userID, result.UserID)

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_GetStoreByID_Success(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	storeID := uint64(1)
	store := &domain.Store{
		ID:          storeID,
		UserID:      1,
		Name:        "Public Store",
		Description: "Public Description",
	}

	// Mock expectations
	mockStoreRepo.On("GetByID", storeID).Return(store, nil)

	// Execute
	result, err := storeUsecase.GetStoreByID(storeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, store.ID, result.ID)
	assert.Equal(t, store.Name, result.Name)

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_GetAllStores_Success(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	page := 1
	limit := 10
	search := "test"
	offset := 0

	stores := []*domain.Store{
		{ID: 1, Name: "Test Store 1"},
		{ID: 2, Name: "Test Store 2"},
	}
	total := int64(2)

	// Mock expectations
	mockStoreRepo.On("GetAll", limit, offset, search).Return(stores, total, nil)

	// Execute
	result, meta, err := storeUsecase.GetAllStores(page, limit, search)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, page, meta.Page)
	assert.Equal(t, limit, meta.Limit)
	assert.Equal(t, total, meta.Total)
	assert.Equal(t, 1, meta.TotalPage)

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_GetAllStores_WithPagination(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	page := 2
	limit := 5
	search := ""
	offset := 5 // (page-1) * limit

	stores := []*domain.Store{
		{ID: 6, Name: "Store 6"},
		{ID: 7, Name: "Store 7"},
	}
	total := int64(12) // Total 12 stores, should have 3 pages

	// Mock expectations
	mockStoreRepo.On("GetAll", limit, offset, search).Return(stores, total, nil)

	// Execute
	result, meta, err := storeUsecase.GetAllStores(page, limit, search)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, page, meta.Page)
	assert.Equal(t, limit, meta.Limit)
	assert.Equal(t, total, meta.Total)
	assert.Equal(t, 3, meta.TotalPage) // ceil(12/5) = 3

	mockStoreRepo.AssertExpectations(t)
}

func TestStoreUsecase_CreateStore_Success(t *testing.T) {
	// Setup
	mockStoreRepo := new(mocks.MockStoreRepository)
	storeUsecase := NewStoreUsecase(mockStoreRepo)

	userID := uint64(1)
	req := &domain.CreateStoreRequest{
		Name:        "New Store",
		Description: "Welcome to New Store",
	}

	// Mock expectations
	mockStoreRepo.On("GetByUserID", userID).Return(nil, gorm.ErrRecordNotFound)
	mockStoreRepo.On("Create", mock.MatchedBy(func(store *domain.Store) bool {
		return store.UserID == userID && store.Name == req.Name
	})).Return(nil)
	mockStoreRepo.On("GetByID", mock.AnythingOfType("uint64")).Return(&domain.Store{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}, nil)

	// Execute
	result, err := storeUsecase.CreateStore(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Description, result.Description)

	mockStoreRepo.AssertExpectations(t)
}