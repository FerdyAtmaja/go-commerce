package usecase

import (
	"errors"
	"math"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"
	"go-commerce/pkg/utils"

	"gorm.io/gorm"
)

type CategoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUsecase) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	// Check if category name already exists
	if _, err := u.categoryRepo.GetByName(req.Name); err == nil {
		return nil, errors.New("category name already exists")
	}

	// Validate parent category exists and is active if provided
	if req.ParentID != nil {
		parent, err := u.categoryRepo.GetByID(*req.ParentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent category not found")
			}
			return nil, errors.New("failed to validate parent category")
		}
		if parent.Status != "active" {
			return nil, errors.New("parent category must be active")
		}
	}

	// Generate unique slug from name
	baseSlug := utils.GenerateSlug(req.Name)
	slug := utils.EnsureUniqueSlug(baseSlug, func(s string) bool {
		_, err := u.categoryRepo.GetBySlug(s)
		return err == nil // true if slug exists
	})

	category := &domain.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
		Slug:     slug,
		Status:   "active",
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	// Update parent's child flags if this is a subcategory
	if req.ParentID != nil {
		go func() {
			u.categoryRepo.UpdateChildFlags(*req.ParentID)
		}()
	}

	return category, nil
}

func (u *CategoryUsecase) GetCategoryByID(id uint64) (*domain.Category, error) {
	category, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	return category, nil
}

func (u *CategoryUsecase) UpdateCategory(id uint64, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	// Get existing category
	existingCategory, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	// Store original name for comparison
	originalName := existingCategory.Name

	// Check if new name already exists (excluding current category)
	if req.Name != originalName {
		if existingCategory, err := u.categoryRepo.GetByName(req.Name); err == nil && existingCategory.ID != id {
			return nil, errors.New("category name already exists")
		}
	}

	// Update category fields
	existingCategory.Name = req.Name
	existingCategory.ParentID = req.ParentID
	// Update slug if name changed
	if req.Name != originalName {
		baseSlug := utils.GenerateSlug(req.Name)
		slug := utils.EnsureUniqueSlug(baseSlug, func(s string) bool {
			existingCategory, err := u.categoryRepo.GetBySlug(s)
			return err == nil && existingCategory.ID != id // true if slug exists and not current category
		})
		existingCategory.Slug = slug
	}

	if err := u.categoryRepo.Update(existingCategory); err != nil {
		return nil, errors.New("failed to update category")
	}

	return existingCategory, nil
}

func (u *CategoryUsecase) DeleteCategory(id uint64) error {
	// Check if category exists
	_, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return errors.New("failed to get category")
	}

	// Check if category has any children (active or inactive)
	children, err := u.categoryRepo.GetChildrenByParentID(id)
	if err != nil {
		return errors.New("failed to check children")
	}
	if len(children) > 0 {
		return errors.New("cannot delete category with child categories")
	}

	// Check if category was ever used by any product (current or historical)
	hasHistorical, err := u.categoryRepo.HasHistoricalProducts(id)
	if err != nil {
		return errors.New("failed to check product history")
	}
	if hasHistorical {
		return errors.New("cannot delete category that has been used by products")
	}

	if err := u.categoryRepo.Delete(id); err != nil {
		return errors.New("failed to delete category")
	}

	return nil
}

func (u *CategoryUsecase) GetAllCategories(page, limit int) ([]*domain.Category, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	categories, total, err := u.categoryRepo.GetAll(limit, offset)
	if err != nil {
		return nil, response.PaginationMeta{}, errors.New("failed to get categories")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	return categories, meta, nil
}

func (u *CategoryUsecase) GetCategoryBySlug(slug string) (*domain.Category, error) {
	category, err := u.categoryRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	return category, nil
}

func (u *CategoryUsecase) GetRootCategories(page, limit int) ([]*domain.Category, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	categories, total, err := u.categoryRepo.GetRootCategories(limit, offset)
	if err != nil {
		return nil, response.PaginationMeta{}, errors.New("failed to get root categories")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	return categories, meta, nil
}

func (u *CategoryUsecase) GetChildrenByParentID(parentID uint64) ([]*domain.Category, error) {
	// Verify parent category exists
	_, err := u.categoryRepo.GetByID(parentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parent category not found")
		}
		return nil, errors.New("failed to verify parent category")
	}

	children, err := u.categoryRepo.GetChildrenByParentID(parentID)
	if err != nil {
		return nil, errors.New("failed to get children categories")
	}

	return children, nil
}

// DeactivateCategory implements business rules for deactivating a category
func (u *CategoryUsecase) DeactivateCategory(categoryID uint64) error {
	category, err := u.categoryRepo.GetByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return errors.New("failed to get category")
	}

	if category.Status == "inactive" {
		return errors.New("category already inactive")
	}

	hasChild, err := u.categoryRepo.HasActiveChildren(categoryID)
	if err != nil {
		return errors.New("failed to check active children")
	}
	if hasChild {
		return errors.New("category has active child categories")
	}

	hasProduct, err := u.categoryRepo.HasActiveProducts(categoryID)
	if err != nil {
		return errors.New("failed to check active products")
	}
	if hasProduct {
		return errors.New("category is used by active products")
	}

	if err := u.categoryRepo.UpdateStatus(categoryID, "inactive"); err != nil {
		return errors.New("failed to deactivate category")
	}

	return nil
}

// ActivateCategory implements business rules for activating a category
func (u *CategoryUsecase) ActivateCategory(categoryID uint64) error {
	category, err := u.categoryRepo.GetByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return errors.New("failed to get category")
	}

	if category.Status == "active" {
		return errors.New("category already active")
	}

	if category.ParentID != nil {
		parentStatus, err := u.categoryRepo.GetParentStatus(categoryID)
		if err != nil {
			return errors.New("failed to check parent status")
		}
		if parentStatus != "active" {
			return errors.New("parent category inactive")
		}
	}

	if err := u.categoryRepo.UpdateStatus(categoryID, "active"); err != nil {
		return errors.New("failed to activate category")
	}

	return nil
}