package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTransactionItemRepository struct {
	mock.Mock
}

func (m *MockTransactionItemRepository) Create(item *domain.TransactionItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockTransactionItemRepository) CreateWithTx(dbTx interface{}, item *domain.TransactionItem) error {
	args := m.Called(dbTx, item)
	return args.Error(0)
}

func (m *MockTransactionItemRepository) GetByTransactionID(transactionID uint64) ([]*domain.TransactionItem, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TransactionItem), args.Error(1)
}