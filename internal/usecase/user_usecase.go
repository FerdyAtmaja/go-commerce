package usecase

import (
	"errors"
	"time"

	"go-commerce/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) GetProfile(userID uint64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
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

func (u *UserUsecase) UpdateProfile(userID uint64, req *domain.UpdateProfileRequest) (*domain.User, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	// Check if phone is being changed and already exists
	if req.Phone != user.Phone {
		if _, err := u.userRepo.GetByPhone(req.Phone); err == nil {
			return nil, errors.New("phone number already registered")
		}
	}

	// Parse date of birth
	var dateOfBirth *time.Time
	if req.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			dateOfBirth = &dob
		}
	}

	// Update user fields
	user.Name = req.Name
	user.Phone = req.Phone
	user.DateOfBirth = dateOfBirth

	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (u *UserUsecase) UpdatePhoto(userID uint64, photoURL string) (*domain.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	user.PhotoURL = photoURL
	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update photo")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (u *UserUsecase) ChangePassword(userID uint64, req *domain.ChangePasswordRequest) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to get user")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	if err := u.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}