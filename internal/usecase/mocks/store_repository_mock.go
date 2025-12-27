package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockStoreRepository struct {
	mock.Mock
}

func (m *MockStoreRepository) Create(store *domain.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *MockStoreRepository) GetByID(id uint64) (*domain.Store, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Store), args.Error(1)
}

func (m *MockStoreRepository) GetByUserID(userID uint64) (*domain.Store, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Store), args.Error(1)
}

func (m *MockStoreRepository) Update(store *domain.Store) error {
	args := m.Called(store)
	return args.Error(0)
}

func (m *MockStoreRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStoreRepository) GetAll(limit, offset int, search string) ([]*domain.Store, int64, error) {
	args := m.Called(limit, offset, search)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Store), args.Get(1).(int64), args.Error(2)
}