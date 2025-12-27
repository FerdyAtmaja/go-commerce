package mocks

import (
	"go-commerce/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockRegionService struct {
	mock.Mock
}

func (m *MockRegionService) GetProvinces() ([]*domain.Province, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Province), args.Error(1)
}

func (m *MockRegionService) GetCitiesByProvince(provinceID string) ([]*domain.City, error) {
	args := m.Called(provinceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.City), args.Error(1)
}

func (m *MockRegionService) ValidateProvinceAndCity(provinceID, cityID string) error {
	args := m.Called(provinceID, cityID)
	return args.Error(0)
}