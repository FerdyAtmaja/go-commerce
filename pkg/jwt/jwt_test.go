package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAccessToken_Success(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", 24, 168)
	
	userID := uint(1)
	email := "test@example.com"
	isAdmin := false

	// Execute
	token, err := jwtManager.GenerateAccessToken(userID, email, isAdmin)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_GenerateRefreshToken_Success(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", 24, 168)
	
	userID := uint(1)

	// Execute
	token, err := jwtManager.GenerateRefreshToken(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_ValidateToken_Success(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", 24, 168)
	
	userID := uint(1)
	email := "test@example.com"
	isAdmin := true

	// Generate token
	token, err := jwtManager.GenerateAccessToken(userID, email, isAdmin)
	assert.NoError(t, err)

	// Execute
	claims, err := jwtManager.ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, isAdmin, claims.IsAdmin)
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	// Setup
	jwtManager := NewJWTManager("test-secret-key", 24, 168)
	
	invalidToken := "invalid.token.here"

	// Execute
	claims, err := jwtManager.ValidateToken(invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	// Setup - Create JWT manager with very short expiration
	jwtManager := NewJWTManager("test-secret-key", 0, 168) // 0 hours = expired immediately
	
	userID := uint(1)
	email := "test@example.com"
	isAdmin := false

	// Generate token (will be expired)
	token, err := jwtManager.GenerateAccessToken(userID, email, isAdmin)
	assert.NoError(t, err)

	// Wait a moment to ensure expiration
	time.Sleep(time.Millisecond * 100)

	// Execute
	claims, err := jwtManager.ValidateToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	// Setup
	jwtManager1 := NewJWTManager("secret-1", 24, 168)
	jwtManager2 := NewJWTManager("secret-2", 24, 168)
	
	userID := uint(1)
	email := "test@example.com"
	isAdmin := false

	// Generate token with first manager
	token, err := jwtManager1.GenerateAccessToken(userID, email, isAdmin)
	assert.NoError(t, err)

	// Try to validate with second manager (different secret)
	claims, err := jwtManager2.ValidateToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}