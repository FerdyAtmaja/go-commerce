package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddressUsecase_CreateAddress_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	userID := uint64(1)
	req := &domain.CreateAddressRequest{
		JudulAlamat:  "Rumah",
		NamaPenerima: "John Doe",
		DetailAlamat: "Jl. Sudirman No. 123",
		NoTelp:       "081234567890",
		KodePos:      "12190",
	}

	// Mock expectations
	mockAddressRepo.On("Create", mock.MatchedBy(func(addr *domain.Address) bool {
		return addr.UserID == userID && addr.JudulAlamat == req.JudulAlamat
	})).Return(nil)

	// Execute
	result, err := addressUsecase.CreateAddress(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, req.JudulAlamat, result.JudulAlamat)

	mockAddressRepo.AssertExpectations(t)
}

func TestAddressUsecase_CreateAddress_InvalidRegion(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	userID := uint64(1)
	req := &domain.CreateAddressRequest{
		JudulAlamat:  "A", // Too short
		NamaPenerima: "John Doe",
		DetailAlamat: "Jl. Sudirman No. 123",
		NoTelp:       "081234567890",
		KodePos:      "12190",
	}

	// Execute
	result, err := addressUsecase.CreateAddress(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "judul_alamat must be between 2 and 255 characters")
}

func TestAddressUsecase_GetAddressByID_Success(t *testing.T) {
	// Setup
	mockAddressRepo := new(mocks.MockAddressRepository)
	mockRegionService := new(mocks.MockRegionService)
	addressUsecase := NewAddressUsecase(mockAddressRepo, mockRegionService)

	addressID := uint64(1)
	userID := uint64(1)
	address := &domain.Address{
		ID:           addressID,
		UserID:       userID,
		JudulAlamat:  "Rumah",
		NamaPenerima: "John Doe",
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

	addressID := uint64(1)
	userID := uint64(2)

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

	addressID := uint64(1)
	userID := uint64(1)

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