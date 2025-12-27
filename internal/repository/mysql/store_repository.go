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

func (r *storeRepository) GetByID(id uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.Preload("User").First(&store, id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetByUserID(userID uint) (*domain.Store, error) {
	var store domain.Store
	err := r.db.Where("user_id = ?", userID).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) Update(store *domain.Store) error {
	return r.db.Save(store).Error
}

func (r *storeRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Store{}, id).Error
}

func (r *storeRepository) GetAll(limit, offset int, search string) ([]*domain.Store, int64, error) {
	var stores []*domain.Store
	var total int64

	query := r.db.Model(&domain.Store{}).Preload("User")
	
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
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