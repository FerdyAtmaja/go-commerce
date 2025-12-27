package usecase

import (
	"errors"
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAddressUsecase_CreateAddress_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	userID := uint(1)
	req := &domain.CreateAddressRequest{
		Name:       "Rumah",
		Detail:     "Jl. Sudirman No. 123",
		Phone:      "081234567890",
		ProvinceID: "31",
		CityID:     "3171",
		PostalCode: "12190",
	}

	// Mock expectations
	mockRegionService.On("ValidateProvinceAndCity", req.ProvinceID, req.CityID).Return(nil)
	mockAddressRepo.On("Create", &domain.Address{
		UserID:     userID,
		Name:       req.Name,
		Detail:     req.Detail,
		Phone:      req.Phone,
		ProvinceID: req.ProvinceID,
		CityID:     req.CityID,
		PostalCode: req.PostalCode,
	}).Return(nil)

	// Execute
	result, err := addressUsecase.CreateAddress(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.ProvinceID, result.ProvinceID)
	assert.Equal(t, req.CityID, result.CityID)

	mockAddressRepo.AssertExpectations(t)
	mockRegionService.AssertExpectations(t)
}

func TestAddressUsecase_CreateAddress_InvalidRegion(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	userID := uint(1)
	req := &domain.CreateAddressRequest{
		Name:       "Rumah",
		Detail:     "Jl. Sudirman No. 123",
		Phone:      "081234567890",
		ProvinceID: "99",
		CityID:     "9999",
		PostalCode: "12190",
	}

	// Mock expectations
	mockRegionService.On("ValidateProvinceAndCity", req.ProvinceID, req.CityID).Return(errors.New("invalid province ID"))

	// Execute
	result, err := addressUsecase.CreateAddress(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid province or city")

	mockRegionService.AssertExpectations(t)
}

func TestAddressUsecase_GetAddressByID_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	addressID := uint(1)
	userID := uint(1)
	address := &domain.Address{
		ID:     addressID,
		UserID: userID,
		Name:   "Rumah",
	}

	// Mock expectations
	mockAddressRepo.On("CheckOwnership", addressID, userID).Return(true)
	mockAddressRepo.On("GetByID", addressID).Return(address, nil)

	// Execute
	result, err := addressUsecase.GetAddressByID(addressID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, address.ID, result.ID)
	assert.Equal(t, address.UserID, result.UserID)

	mockAddressRepo.AssertExpectations(t)
}

func TestAddressUsecase_GetAddressByID_AccessDenied(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	addressID := uint(1)
	userID := uint(2)

	// Mock expectations
	mockAddressRepo.On("CheckOwnership", addressID, userID).Return(false)

	// Execute
	result, err := addressUsecase.GetAddressByID(addressID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "address not found or access denied")

	mockAddressRepo.AssertExpectations(t)
}

func TestAddressUsecase_DeleteAddress_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	addressID := uint(1)
	userID := uint(1)

	// Mock expectations
	mockAddressRepo.On("CheckOwnership", addressID, userID).Return(true)
	mockAddressRepo.On("Delete", addressID).Return(nil)

	// Execute
	err := addressUsecase.DeleteAddress(addressID, userID)

	// Assert
	assert.NoError(t, err)

	mockAddressRepo.AssertExpectations(t)
}

func TestAddressUsecase_GetProvinces_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	provinces := []*domain.Province{
		{ID: "31", Name: "DKI JAKARTA"},
		{ID: "32", Name: "JAWA BARAT"},
	}

	// Mock expectations
	mockRegionService.On("GetProvinces").Return(provinces, nil)

	// Execute
	result, err := addressUsecase.GetProvinces()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "31", result[0].ID)
	assert.Equal(t, "DKI JAKARTA", result[0].Name)

	mockRegionService.AssertExpectations(t)
}