package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockProductLogRepository struct {
	mock.Mock
}

func (m *MockProductLogRepository) Create(log *domain.ProductLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockProductLogRepository) CreateAsync(log *domain.ProductLog) {
	m.Called(log)
}

func (m *MockProductLogRepository) GetByProductID(productID uint64) ([]*domain.ProductLog, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ProductLog), args.Error(1)
}

func (m *MockProductLogRepository) GetByID(id uint64) (*domain.ProductLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductLog), args.Error(1)
}