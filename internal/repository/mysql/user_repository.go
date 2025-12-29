package mysql

import (
	"go-commerce/internal/domain"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint64) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("notelp = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepository) UpdateLastLogin(userID uint64, lastLogin time.Time) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("last_login_at", lastLogin).Error
}

func (r *userRepository) UpdateProfile(userID uint64, updates map[string]interface{}) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error
}