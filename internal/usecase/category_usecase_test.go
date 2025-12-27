package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryUsecase_CreateCategory_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	req := &domain.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Electronic products",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByName", req.Name).Return(nil, gorm.ErrRecordNotFound) // Name not exists
	mockCategoryRepo.On("Create", &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}).Return(nil)

	// Execute
	result, err := categoryUsecase.CreateCategory(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Description, result.Description)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_CreateCategory_NameAlreadyExists(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	req := &domain.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Electronic products",
	}

	existingCategory := &domain.Category{
		ID:   1,
		Name: req.Name,
	}

	// Mock expectations
	mockCategoryRepo.On("GetByName", req.Name).Return(existingCategory, nil) // Name exists

	// Execute
	result, err := categoryUsecase.CreateCategory(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category name already exists")

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_GetCategoryByID_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint(1)
	category := &domain.Category{
		ID:          categoryID,
		Name:        "Electronics",
		Description: "Electronic products",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)

	// Execute
	result, err := categoryUsecase.GetCategoryByID(categoryID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, category.ID, result.ID)
	assert.Equal(t, category.Name, result.Name)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_GetCategoryByID_NotFound(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint(999)

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := categoryUsecase.GetCategoryByID(categoryID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category not found")

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_UpdateCategory_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint(1)
	existingCategory := &domain.Category{
		ID:          categoryID,
		Name:        "Electronics",
		Description: "Old description",
	}

	req := &domain.UpdateCategoryRequest{
		Name:        "Updated Electronics",
		Description: "Updated description",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(existingCategory, nil)
	mockCategoryRepo.On("GetByName", req.Name).Return(nil, gorm.ErrRecordNotFound) // New name not exists
	mockCategoryRepo.On("Update", existingCategory).Return(nil)

	// Execute
	result, err := categoryUsecase.UpdateCategory(categoryID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Description, result.Description)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_UpdateCategory_NameAlreadyExists(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint(1)
	existingCategory := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
	}

	anotherCategory := &domain.Category{
		ID:   2,
		Name: "Clothing",
	}

	req := &domain.UpdateCategoryRequest{
		Name: "Clothing", // This name belongs to another category
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(existingCategory, nil)
	mockCategoryRepo.On("GetByName", req.Name).Return(anotherCategory, nil) // Name exists for different category

	// Execute
	result, err := categoryUsecase.UpdateCategory(categoryID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category name already exists")

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_DeleteCategory_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint(1)
	category := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockCategoryRepo.On("Delete", categoryID).Return(nil)

	// Execute
	err := categoryUsecase.DeleteCategory(categoryID)

	// Assert
	assert.NoError(t, err)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_GetAllCategories_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	page := 1
	limit := 10
	offset := 0

	categories := []*domain.Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Clothing"},
	}
	total := int64(2)

	// Mock expectations
	mockCategoryRepo.On("GetAll", limit, offset).Return(categories, total, nil)

	// Execute
	result, meta, err := categoryUsecase.GetAllCategories(page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, page, meta.Page)
	assert.Equal(t, limit, meta.Limit)
	assert.Equal(t, total, meta.Total)
	assert.Equal(t, 1, meta.TotalPage)

	mockCategoryRepo.AssertExpectations(t)
}