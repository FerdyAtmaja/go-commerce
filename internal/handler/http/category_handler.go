package http

import (
	"strconv"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
	validator       *validator.Validate
}

func NewCategoryHandler(categoryUsecase *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
		validator:       validator.New(),
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (admin only). For root category, omit parent_id or set to null. For subcategory, provide parent_id.
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} response.Response{data=domain.Category} "Category created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req domain.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.ParentID != nil && *req.ParentID == 0 {
		req.ParentID = nil
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	category, err := h.categoryUsecase.CreateCategory(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Category created successfully", category)
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Get all product categories with pagination (public endpoint)
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Category} "Categories retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	categories, meta, err := h.categoryUsecase.GetAllCategories(page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Categories retrieved successfully", categories, meta)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a single category by its ID (public endpoint)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response{data=domain.Category} "Category retrieved successfully"
// @Failure 400 {object} response.Response "Invalid category ID"
// @Failure 404 {object} response.Response "Category not found"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	category, err := h.categoryUsecase.GetCategoryByID(id)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Category retrieved successfully", category)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body domain.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} response.Response{data=domain.Category} "Category updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	var req domain.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	category, err := h.categoryUsecase.UpdateCategory(id, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Category updated successfully", category)
}

// DeactivateCategory godoc
// @Summary Deactivate a category
// @Description Deactivate a category following business rules: cannot deactivate if has active children or used by active products (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response "Category deactivated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - category has active children or used by active products"
// @Router /categories/{id}/deactivate [put]
func (h *CategoryHandler) DeactivateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	if err := h.categoryUsecase.DeactivateCategory(id); err != nil {
		return response.Conflict(c, err.Error())
	}

	return response.Success(c, "Category deactivated successfully", nil)
}

// ActivateCategory godoc
// @Summary Activate a category
// @Description Activate a category following business rules: cannot activate if parent is inactive (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response "Category activated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - parent category is inactive"
// @Router /categories/{id}/activate [put]
func (h *CategoryHandler) ActivateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	if err := h.categoryUsecase.ActivateCategory(id); err != nil {
		return response.Conflict(c, err.Error())
	}

	return response.Success(c, "Category activated successfully", nil)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Physical delete of category (admin only). Only allowed if never used by any product in history
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response "Category deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only or category has been used"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	if err := h.categoryUsecase.DeleteCategory(id); err != nil {
		return response.Forbidden(c, err.Error())
	}

	return response.Success(c, "Category deleted successfully", nil)
}