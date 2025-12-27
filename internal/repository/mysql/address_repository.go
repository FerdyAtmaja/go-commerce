package mysql

import (
	"go-commerce/internal/domain"

	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) domain.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *domain.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) GetByID(id uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.Preload("User").First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetByUserID(userID uint, limit, offset int) ([]*domain.Address, int64, error) {
	var addresses []*domain.Address
	var total int64

	query := r.db.Model(&domain.Address{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Limit(limit).Offset(offset).Find(&addresses).Error; err != nil {
		return nil, 0, err
	}

	return addresses, total, nil
}

func (r *addressRepository) Update(address *domain.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

func (r *addressRepository) CheckOwnership(addressID, userID uint) bool {
	var count int64
	r.db.Model(&domain.Address{}).Where("id = ? AND user_id = ?", addressID, userID).Count(&count)
	return count > 0
}