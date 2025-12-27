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

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req domain.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
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

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	categories, meta, err := h.categoryUsecase.GetAllCategories(page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Categories retrieved successfully", categories, meta)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	category, err := h.categoryUsecase.GetCategoryByID(uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
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

	category, err := h.categoryUsecase.UpdateCategory(uint(id), &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Category updated successfully", category)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	if err := h.categoryUsecase.DeleteCategory(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Category deleted successfully", nil)
}