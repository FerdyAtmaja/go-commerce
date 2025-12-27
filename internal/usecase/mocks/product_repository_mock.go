package mocks

import (
	"go-commerce/internal/domain"
	"github.com/stretchr/testify/mock"
)

// Alias untuk compatibility dengan mock yang sudah ada
type StoreRepositoryMock = MockStoreRepository
type CategoryRepositoryMock = MockCategoryRepository

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetByID(id uint64) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetBySlug(slug string) (*domain.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByTokoID(tokoID uint64, limit, offset int, search string) ([]*domain.Product, int64, error) {
	args := m.Called(tokoID, limit, offset, search)
	return args.Get(0).([]*domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) GetAll(limit, offset int, search, categoryID string) ([]*domain.Product, int64, error) {
	args := m.Called(limit, offset, search, categoryID)
	return args.Get(0).([]*domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) CheckOwnership(productID, tokoID uint64) error {
	args := m.Called(productID, tokoID)
	return args.Error(0)
}

type PhotoProdukRepositoryMock struct {
	mock.Mock
}

func (m *PhotoProdukRepositoryMock) Create(photo *domain.PhotoProduk) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *PhotoProdukRepositoryMock) GetByProductID(productID uint64) ([]*domain.PhotoProduk, error) {
	args := m.Called(productID)
	return args.Get(0).([]*domain.PhotoProduk), args.Error(1)
}

func (m *PhotoProdukRepositoryMock) Update(photo *domain.PhotoProduk) error {
	args := m.Called(photo)
	return args.Error(0)
}

func (m *PhotoProdukRepositoryMock) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *PhotoProdukRepositoryMock) SetPrimary(productID, photoID uint64) error {
	args := m.Called(productID, photoID)
	return args.Error(0)
}