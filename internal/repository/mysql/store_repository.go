package mysql

import (
	"go-commerce/internal/domain"

	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) domain.StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(store *domain.Store) error {
	return r.db.Create(store).Error
}

func (r *storeRepository) GetByID(id uint64) (*domain.Store, error) {
	var store domain.Store
	err := r.db.Preload("User").First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetByUserID(userID uint64) (*domain.Store, error) {
	var store domain.Store
	err := r.db.Where("id_user = ?", userID).Preload("User").First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) Update(store *domain.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Store{}, id).Error
}

func (r *storeRepository) GetAll(limit, offset int, search string) ([]*domain.Store, int64, error) {
	var stores []*domain.Store
	var total int64

	query := r.db.Model(&domain.Store{}).Preload("User")
	
	if search != "" {
		query = query.Where("nama_toko LIKE ?", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Limit(limit).Offset(offset).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// GetActiveStores returns only active stores for public listing
func (r *storeRepository) GetActiveStores(limit, offset int, search string) ([]*domain.Store, int64, error) {
	var stores []*domain.Store
	var total int64

	query := r.db.Model(&domain.Store{}).Preload("User").Where("status = ?", "active")
	
	if search != "" {
		query = query.Where("nama_toko LIKE ?", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Limit(limit).Offset(offset).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}
// GetPendingStores returns stores waiting for admin approval
func (r *storeRepository) GetPendingStores(limit, offset int, search string) ([]*domain.Store, int64, error) {
	var stores []*domain.Store
	var total int64

	query := r.db.Model(&domain.Store{}).Preload("User").Where("status = ?", "pending")
	
	if search != "" {
		query = query.Where("nama_toko LIKE ?", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := query.Limit(limit).Offset(offset).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}
// GetActiveStoreByUserID returns active store for a user
func (r *storeRepository) GetActiveStoreByUserID(userID uint64) (*domain.Store, error) {
	var store domain.Store
	err := r.db.Where("id_user = ? AND status IN (?)", userID, []string{"active", "pending"}).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}