package usecase

import (
	"errors"
	"math"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"

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

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	return category, nil
}

func (u *CategoryUsecase) GetCategoryByID(id uint) (*domain.Category, error) {
	category, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	return category, nil
}

func (u *CategoryUsecase) UpdateCategory(id uint, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	// Get existing category
	category, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, errors.New("failed to get category")
	}

	// Check if new name already exists (excluding current category)
	if req.Name != category.Name {
		if existingCategory, err := u.categoryRepo.GetByName(req.Name); err == nil && existingCategory.ID != id {
			return nil, errors.New("category name already exists")
		}
	}

	// Update category fields
	category.Name = req.Name
	category.Description = req.Description

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, errors.New("failed to update category")
	}

	return category, nil
}

func (u *CategoryUsecase) DeleteCategory(id uint) error {
	// Check if category exists
	_, err := u.categoryRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return errors.New("failed to get category")
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
	if limit < 1 {
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