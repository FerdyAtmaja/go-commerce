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

func (r *addressRepository) GetByID(id uint64) (*domain.Address, error) {
	var address domain.Address
	err := r.db.Preload("User").First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetByUserID(userID uint64, limit, offset int) ([]*domain.Address, int64, error) {
	var addresses []*domain.Address
	var total int64

	query := r.db.Model(&domain.Address{}).Where("id_user = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination, ordered by is_default DESC, created_at DESC
	if err := query.Order("is_default DESC, created_at DESC").Limit(limit).Offset(offset).Find(&addresses).Error; err != nil {
		return nil, 0, err
	}

	return addresses, total, nil
}

func (r *addressRepository) GetDefaultByUserID(userID uint64) (*domain.Address, error) {
	var address domain.Address
	err := r.db.Where("id_user = ? AND is_default = ?", userID, true).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Update(address *domain.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

func (r *addressRepository) SetDefault(addressID, userID uint64) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all default addresses for user
	if err := tx.Model(&domain.Address{}).Where("id_user = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set new default address
	if err := tx.Model(&domain.Address{}).Where("id = ? AND id_user = ?", addressID, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *addressRepository) CheckOwnership(addressID, userID uint64) bool {
	var count int64
	r.db.Model(&domain.Address{}).Where("id = ? AND id_user = ?", addressID, userID).Count(&count)
	return count > 0
}