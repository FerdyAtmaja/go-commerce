package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(id uint64) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(name string) (*domain.Category, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetAll(limit, offset int) ([]*domain.Category, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetRootCategories(limit, offset int) ([]*domain.Category, int64, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Category), args.Get(1).(int64), args.Error(2)
}

func (m *MockCategoryRepository) GetChildrenByParentID(parentID uint64) ([]*domain.Category, error) {
	args := m.Called(parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) HasActiveChildren(categoryID uint64) (bool, error) {
	args := m.Called(categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryRepository) HasActiveProducts(categoryID uint64) (bool, error) {
	args := m.Called(categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryRepository) HasHistoricalProducts(categoryID uint64) (bool, error) {
	args := m.Called(categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryRepository) UpdateStatus(categoryID uint64, status string) error {
	args := m.Called(categoryID, status)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetParentStatus(categoryID uint64) (string, error) {
	args := m.Called(categoryID)
	return args.String(0), args.Error(1)
}

func (m *MockCategoryRepository) UpdateHasActiveProduct(categoryID uint64) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

func (m *MockCategoryRepository) UpdateChildFlags(categoryID uint64) error {
	args := m.Called(categoryID)
	return args.Error(0)
}