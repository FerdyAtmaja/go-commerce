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
// @Param request body domain.CreateCategoryRequest true "Category creation request. For root category: {\"name\": \"Electronics\", \"parent_id\": 0} or omit parent_id. For subcategory: {\"name\": \"Smartphones\", \"parent_id\": 1}"
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

	// Handle special case: if parent_id is 0, treat as null (root category)
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

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an existing category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response "Category deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	if err := h.categoryUsecase.DeleteCategory(id); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Category deleted successfully", nil)
}

// GetCategoryBySlug godoc
// @Summary Get category by slug
// @Description Get a single category by its slug (public endpoint)
// @Tags Categories
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} response.Response{data=domain.Category} "Category retrieved successfully"
// @Failure 404 {object} response.Response "Category not found"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.BadRequest(c, "Category slug is required")
	}

	category, err := h.categoryUsecase.GetCategoryBySlug(slug)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Category retrieved successfully", category)
}

// GetRootCategories godoc
// @Summary Get root categories
// @Description Get all root categories (categories without parent) with pagination (public endpoint)
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Category} "Root categories retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	categories, meta, err := h.categoryUsecase.GetRootCategories(page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Root categories retrieved successfully", categories, meta)
}

// GetChildrenByParentID godoc
// @Summary Get children categories
// @Description Get all children categories of a parent category (public endpoint)
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Parent Category ID"
// @Success 200 {object} response.Response{data=[]domain.Category} "Children categories retrieved successfully"
// @Failure 400 {object} response.Response "Invalid parent category ID"
// @Failure 404 {object} response.Response "Parent category not found"
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetChildrenByParentID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid parent category ID")
	}

	children, err := h.categoryUsecase.GetChildrenByParentID(id)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Children categories retrieved successfully", children)
}