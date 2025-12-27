package usecase

import (
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"
	"go-commerce/pkg/jwt"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthUsecase_Login_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.MockStoreRepository)
	jwtManager := jwt.NewJWTManager("test-secret", 24, 168)
	
	authUsecase := &AuthUsecase{
		userRepo:   mockUserRepo,
		storeRepo:  mockStoreRepo,
		jwtManager: jwtManager,
		db:         nil, // Not needed for login test
	}

	// Test data
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	req := &domain.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Mock expectations
	mockUserRepo.On("GetByEmail", email).Return(user, nil)

	// Execute
	result, err := authUsecase.Login(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, user.ID, result.User.ID)
	assert.Equal(t, user.Email, result.User.Email)
	assert.Empty(t, result.User.Password) // Password should be removed

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidEmail(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.MockStoreRepository)
	jwtManager := jwt.NewJWTManager("test-secret", 24, 168)
	
	authUsecase := &AuthUsecase{
		userRepo:   mockUserRepo,
		storeRepo:  mockStoreRepo,
		jwtManager: jwtManager,
		db:         nil,
	}

	req := &domain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockUserRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := authUsecase.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid email or password")

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.MockStoreRepository)
	jwtManager := jwt.NewJWTManager("test-secret", 24, 168)
	
	authUsecase := &AuthUsecase{
		userRepo:   mockUserRepo,
		storeRepo:  mockStoreRepo,
		jwtManager: jwtManager,
		db:         nil,
	}

	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	
	user := &domain.User{
		ID:       1,
		Email:    email,
		Password: string(hashedPassword),
	}

	req := &domain.LoginRequest{
		Email:    email,
		Password: wrongPassword,
	}

	// Mock expectations
	mockUserRepo.On("GetByEmail", email).Return(user, nil)

	// Execute
	result, err := authUsecase.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid email or password")

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserByID_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.MockStoreRepository)
	jwtManager := jwt.NewJWTManager("test-secret", 24, 168)
	
	authUsecase := &AuthUsecase{
		userRepo:   mockUserRepo,
		storeRepo:  mockStoreRepo,
		jwtManager: jwtManager,
		db:         nil,
	}

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
	result, err := authUsecase.GetUserByID(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.Empty(t, result.Password) // Password should be removed

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserByID_NotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(mocks.MockUserRepository)
	mockStoreRepo := new(mocks.MockStoreRepository)
	jwtManager := jwt.NewJWTManager("test-secret", 24, 168)
	
	authUsecase := &AuthUsecase{
		userRepo:   mockUserRepo,
		storeRepo:  mockStoreRepo,
		jwtManager: jwtManager,
		db:         nil,
	}

	userID := uint(999)

	// Mock expectations
	mockUserRepo.On("GetByID", userID).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := authUsecase.GetUserByID(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
}