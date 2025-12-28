package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCategoryUsecase_CreateCategory_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	req := &domain.CreateCategoryRequest{
		Name: "Electronics",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByName", req.Name).Return(nil, gorm.ErrRecordNotFound) // Name not exists
	mockCategoryRepo.On("Create", mock.MatchedBy(func(category *domain.Category) bool {
		return category.Name == req.Name && category.Slug == "electronics"
	})).Return(nil)

	// Execute
	result, err := categoryUsecase.CreateCategory(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, "electronics", result.Slug)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_CreateCategory_NameAlreadyExists(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	req := &domain.CreateCategoryRequest{
		Name: "Electronics",
	}

	existingCategory := &domain.Category{
		ID:   1,
		Name: req.Name,
		Slug: "electronics",
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

	categoryID := uint64(1)
	category := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
		Slug: "electronics",
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

	categoryID := uint64(999)

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

	categoryID := uint64(1)
	existingCategory := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
		Slug: "electronics",
	}

	req := &domain.UpdateCategoryRequest{
		Name: "Updated Electronics",
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", categoryID).Return(existingCategory, nil)
	mockCategoryRepo.On("GetByName", req.Name).Return(nil, gorm.ErrRecordNotFound) // New name not exists
	mockCategoryRepo.On("Update", mock.MatchedBy(func(category *domain.Category) bool {
		return category.Name == req.Name && category.Slug == "updated-electronics"
	})).Return(nil)

	// Execute
	result, err := categoryUsecase.UpdateCategory(categoryID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, "updated-electronics", result.Slug)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_UpdateCategory_NameAlreadyExists(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	categoryID := uint64(1)
	existingCategory := &domain.Category{
		ID:   categoryID,
		Name: "Electronics",
		Slug: "electronics",
	}

	anotherCategory := &domain.Category{
		ID:   2,
		Name: "Clothing",
		Slug: "clothing",
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

	categoryID := uint64(1)
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

func TestCategoryUsecase_GetCategoryBySlug_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	slug := "electronics"
	category := &domain.Category{
		ID:   1,
		Name: "Electronics",
		Slug: slug,
	}

	// Mock expectations
	mockCategoryRepo.On("GetBySlug", slug).Return(category, nil)

	// Execute
	result, err := categoryUsecase.GetCategoryBySlug(slug)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, category.Slug, result.Slug)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_GetRootCategories_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	page := 1
	limit := 10
	offset := 0

	rootCategories := []*domain.Category{
		{ID: 1, Name: "Electronics", ParentID: nil},
		{ID: 2, Name: "Clothing", ParentID: nil},
	}
	total := int64(2)

	// Mock expectations
	mockCategoryRepo.On("GetRootCategories", limit, offset).Return(rootCategories, total, nil)

	// Execute
	result, meta, err := categoryUsecase.GetRootCategories(page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, page, meta.Page)
	assert.Equal(t, total, meta.Total)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryUsecase_GetChildrenByParentID_Success(t *testing.T) {
	// Setup
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	categoryUsecase := NewCategoryUsecase(mockCategoryRepo)

	parentID := uint64(1)
	parentCategory := &domain.Category{
		ID:   parentID,
		Name: "Electronics",
	}

	children := []*domain.Category{
		{ID: 2, Name: "Smartphones", ParentID: &parentID},
		{ID: 3, Name: "Laptops", ParentID: &parentID},
	}

	// Mock expectations
	mockCategoryRepo.On("GetByID", parentID).Return(parentCategory, nil)
	mockCategoryRepo.On("GetChildrenByParentID", parentID).Return(children, nil)

	// Execute
	result, err := categoryUsecase.GetChildrenByParentID(parentID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, parentID, *result[0].ParentID)

	mockCategoryRepo.AssertExpectations(t)
}