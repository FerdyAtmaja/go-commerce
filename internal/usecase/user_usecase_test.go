package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserUsecase_GetProfile_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(1)
	user := &domain.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(user, nil)

	// Execute
	result, err := userUsecase.GetProfile(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Empty(t, result.Password) // Password should be removed

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_GetProfile_NotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(999)

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := userUsecase.GetProfile(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateProfile_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(1)
	existingUser := &domain.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "test@example.com",
		Phone: "081234567890",
	}

	req := &domain.UpdateProfileRequest{
		Name:        "New Name",
		Phone:       "081234567891",
		DateOfBirth: "1990-01-01",
	}

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(existingUser, nil)
	mockUserRepo.On("GetByPhone", req.Phone).Return(nil, gorm.ErrRecordNotFound) // Phone not exists
	mockUserRepo.On("Update", existingUser).Return(nil)

	// Execute
	result, err := userUsecase.UpdateProfile(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Phone, result.Phone)
	assert.Empty(t, result.Password)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateProfile_PhoneAlreadyExists(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(1)
	existingUser := &domain.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
		Phone: "081234567890",
	}

	anotherUser := &domain.User{
		ID:    2,
		Phone: "081234567891",
	}

	req := &domain.UpdateProfileRequest{
		Name:  "Updated Name",
		Phone: "081234567891", // This phone belongs to another user
	}

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(existingUser, nil)
	mockUserRepo.On("GetByPhone", req.Phone).Return(anotherUser, nil) // Phone exists

	// Execute
	result, err := userUsecase.UpdateProfile(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "phone number already registered")

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_ChangePassword_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(1)
	currentPassword := "oldpassword"
	newPassword := "newpassword"
	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       userID,
		Password: string(hashedCurrentPassword),
	}

	req := &domain.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(user, nil)
	mockUserRepo.On("Update", user).Return(nil)

	// Execute
	err := userUsecase.ChangePassword(userID, req)

	// Assert
	assert.NoError(t, err)

	// Verify new password is hashed correctly
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

func TestUserUsecase_ChangePassword_WrongCurrentPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	userUsecase := NewUserUsecase(mockUserRepo)

	userID := uint(1)
	correctPassword := "correctpassword"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	user := &domain.User{
		ID:       userID,
		Password: string(hashedPassword),
	}

	req := &domain.ChangePasswordRequest{
		CurrentPassword: wrongPassword,
		NewPassword:     "newpassword",
	}

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(user, nil)

	// Execute
	err := userUsecase.ChangePassword(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "current password is incorrect")

	mockUserRepo.AssertExpectations(t)
}