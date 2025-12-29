package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(id uint64) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(userID uint64, limit, offset int) ([]*domain.Transaction, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) GetByStoreID(storeID uint64, limit, offset int) ([]*domain.Transaction, int64, error) {
	args := m.Called(storeID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) Update(tx *domain.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) BeginTx() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func (m *MockTransactionRepository) CommitTx(tx interface{}) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) RollbackTx(tx interface{}) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateWithTx(dbTx interface{}, tx *domain.Transaction) error {
	args := m.Called(dbTx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByStatus(status string, limit, offset int) ([]*domain.Transaction, int64, error) {
	args := m.Called(status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Get(1).(int64), args.Error(2)
}