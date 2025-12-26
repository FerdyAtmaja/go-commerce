package usecase

import (
	"errors"
	"time"

	"go-commerce/internal/domain"
	"go-commerce/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	userRepo   domain.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtManager *jwt.JWTManager) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (u *AuthUsecase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if email already exists
	if _, err := u.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists
	if _, err := u.userRepo.GetByPhone(req.Phone); err == nil {
		return nil, errors.New("phone number already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			dateOfBirth = &dob
		}
	}

	// Create user
	user := &domain.User{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    string(hashedPassword),
		DateOfBirth: dateOfBirth,
		IsAdmin:     false,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Remove password from response
	user.Password = ""

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (u *AuthUsecase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, errors.New("failed to get user")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := u.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Remove password from response
	user.Password = ""

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (u *AuthUsecase) GetUserByID(id uint) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}