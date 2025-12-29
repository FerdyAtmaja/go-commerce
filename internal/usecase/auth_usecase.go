package usecase

import (
	"errors"
	"log"
	"time"

	"go-commerce/internal/domain"
	"go-commerce/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	userRepo   domain.UserRepository
	storeRepo  domain.StoreRepository
	jwtManager *jwt.JWTManager
	db         *gorm.DB
}

func NewAuthUsecase(userRepo domain.UserRepository, storeRepo domain.StoreRepository, jwtManager *jwt.JWTManager, db *gorm.DB) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		storeRepo:  storeRepo,
		jwtManager: jwtManager,
		db:         db,
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

	// Start database transaction
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create user
	user := &domain.User{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    string(hashedPassword),
		DateOfBirth: dateOfBirth,
		IsAdmin:     false,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("failed to create user")
	}

	// Auto create store
	store := &domain.Store{
		UserID:      user.ID,
		Name:        req.Name + "'s Store",
		Description: "Welcome to " + req.Name + "'s Store",
	}

	if err := tx.Create(store).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating store: %v", err)
		return nil, errors.New("failed to create store")
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to complete registration")
	}

	// Send welcome email async
	go func() {
		// Simulate sending welcome email
		time.Sleep(100 * time.Millisecond)
		// log.Printf("Welcome email sent to %s", user.Email)
	}()

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

	// Update last login async
	go func() {
		now := time.Now()
		u.userRepo.UpdateLastLogin(user.ID, now)
	}()

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

func (u *AuthUsecase) GetUserByID(id uint64) (*domain.User, error) {
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
