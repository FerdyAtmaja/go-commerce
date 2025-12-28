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
	err := r.db.Preload("Parent").Preload("Children").First(&category, id).Error
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

func (r *categoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Parent").Preload("Children").Where("slug = ?", slug).First(&category).Error
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

	// Get data with pagination and preload relations
	if err := r.db.Preload("Parent").Preload("Children").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) GetRootCategories(limit, offset int) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	// Count total root categories (parent_id IS NULL)
	if err := r.db.Model(&domain.Category{}).Where("parent_id IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get root categories with their children
	if err := r.db.Preload("Children").Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) GetChildrenByParentID(parentID uint64) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}