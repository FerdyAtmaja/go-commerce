package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) Create(address *domain.Address) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *MockAddressRepository) GetByID(id uint64) (*domain.Address, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Address), args.Error(1)
}

func (m *MockAddressRepository) GetByUserID(userID uint64, limit, offset int) ([]*domain.Address, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Address), args.Get(1).(int64), args.Error(2)
}

func (m *MockAddressRepository) Update(address *domain.Address) error {
	args := m.Called(address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAddressRepository) CheckOwnership(addressID, userID uint64) bool {
	args := m.Called(addressID, userID)
	return args.Bool(0)
}