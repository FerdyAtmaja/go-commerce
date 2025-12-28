package http

import (
	"strconv"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/middleware"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	storeUsecase *usecase.StoreUsecase
	validator    *validator.Validate
}

func NewStoreHandler(storeUsecase *usecase.StoreUsecase) *StoreHandler {
	return &StoreHandler{
		storeUsecase: storeUsecase,
		validator:    validator.New(),
	}
}

// GetMyStore godoc
// @Summary Get current user's store
// @Description Get store information for the authenticated user
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.Store} "Store retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Store not found"
// @Router /stores/my [get]
func (h *StoreHandler) GetMyStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	store, err := h.storeUsecase.GetMyStore(userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Store retrieved successfully", store)
}

// UpdateMyStore godoc
// @Summary Update current user's store
// @Description Update store information for the authenticated user
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateStoreRequest true "Store update request"
// @Success 200 {object} response.Response{data=domain.Store} "Store updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /stores/my [put]
func (h *StoreHandler) UpdateMyStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req domain.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	store, err := h.storeUsecase.UpdateMyStore(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store updated successfully", store)
}

// GetAllStores godoc
// @Summary Get all stores
// @Description Get all stores with pagination and search (public endpoint)
// @Tags Stores
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by store name"
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Store} "Stores retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /stores [get]
func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	stores, meta, err := h.storeUsecase.GetAllStores(page, limit, search)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Stores retrieved successfully", stores, meta)
}

// GetStoreByID godoc
// @Summary Get store by ID
// @Description Get a single store by its ID (public endpoint)
// @Tags Stores
// @Accept json
// @Produce json
// @Param id path int true "Store ID"
// @Success 200 {object} response.Response{data=domain.Store} "Store retrieved successfully"
// @Failure 400 {object} response.Response "Invalid store ID"
// @Failure 404 {object} response.Response "Store not found"
// @Router /stores/{id} [get]
func (h *StoreHandler) GetStoreByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid store ID")
	}

	store, err := h.storeUsecase.GetStoreByID(id)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Store retrieved successfully", store)
}