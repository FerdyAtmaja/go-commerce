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

func (r *categoryRepository) HasActiveChildren(categoryID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Category{}).Where("parent_id = ? AND status = 'active'", categoryID).Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) HasActiveProducts(categoryID uint64) (bool, error) {
	var count int64
	err := r.db.Table("produk").Where("id_category = ? AND status = 'active'", categoryID).Count(&count).Error
	return count > 0, err
}

// HasHistoricalProducts checks if category was ever used by any product (including in log_produk)
func (r *categoryRepository) HasHistoricalProducts(categoryID uint64) (bool, error) {
	// Check current products
	var currentCount int64
	err := r.db.Table("produk").Where("id_category = ?", categoryID).Count(&currentCount).Error
	if err != nil {
		return false, err
	}
	if currentCount > 0 {
		return true, nil
	}
	
	// Check historical products in log
	var logCount int64
	err = r.db.Table("log_produk").Where("id_category = ?", categoryID).Count(&logCount).Error
	return logCount > 0, err
}

func (r *categoryRepository) UpdateStatus(categoryID uint64, status string) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", categoryID).Update("status", status).Error
}

func (r *categoryRepository) GetParentStatus(categoryID uint64) (string, error) {
	var category domain.Category
	err := r.db.Select("parent_id").First(&category, categoryID).Error
	if err != nil {
		return "", err
	}
	
	if category.ParentID == nil {
		return "active", nil // Root category, assume active
	}
	
	var parent domain.Category
	err = r.db.Select("status").First(&parent, *category.ParentID).Error
	if err != nil {
		return "", err
	}
	
	return parent.Status, nil
}

// UpdateHasActiveProduct updates the has_active_product flag for a category
func (r *categoryRepository) UpdateHasActiveProduct(categoryID uint64) error {
	hasActive, err := r.HasActiveProducts(categoryID)
	if err != nil {
		return err
	}
	
	return r.db.Model(&domain.Category{}).Where("id = ?", categoryID).Update("has_active_product", hasActive).Error
}

// UpdateChildFlags updates has_child and is_leaf flags for a category
func (r *categoryRepository) UpdateChildFlags(categoryID uint64) error {
	hasChild, err := r.HasActiveChildren(categoryID)
	if err != nil {
		return err
	}
	
	return r.db.Model(&domain.Category{}).Where("id = ?", categoryID).Updates(map[string]interface{}{
		"has_child": hasChild,
		"is_leaf":   !hasChild,
	}).Error
}