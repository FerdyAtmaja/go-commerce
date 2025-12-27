package mysql

import (
	"go-commerce/internal/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint64) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("nama_category = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *categoryRepository) GetAll(limit, offset int) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	// Count total
	if err := r.db.Model(&domain.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	if err := r.db.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}