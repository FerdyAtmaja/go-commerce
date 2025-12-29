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
// @Summary Get current user's store (Seller only)
// @Description Get store information for the authenticated user. Only store owners can access their store.
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
// @Summary Update current user's store (Seller only)
// @Description Update store information for the authenticated user. Only store owners can update their store.
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
// @Summary Get all active stores (Public)
// @Description Get all active stores with pagination and search. Only shows active stores to public.
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

	stores, meta, err := h.storeUsecase.GetActiveStores(page, limit, search)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Stores retrieved successfully", stores, meta)
}

// GetStoreByID godoc
// @Summary Get store by ID (Public)
// @Description Get a single active store by its ID. Only shows active stores to public.
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

	store, err := h.storeUsecase.GetStorePublic(id)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Store retrieved successfully", store)
}

// CreateStore godoc
// @Summary Create a new store (Authenticated User)
// @Description Create a new store for the authenticated user. Requires authentication.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateStoreRequest true "Store creation request"
// @Success 201 {object} response.Response{data=domain.Store} "Store created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /stores [post]
func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req domain.CreateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	store, err := h.storeUsecase.CreateStore(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Store created successfully", store)
}
// ActivateStore godoc
// @Summary Activate store (Seller only)
// @Description Activate the authenticated user's store. Store must not be suspended and profile must be complete.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Store activated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - store suspended"
// @Failure 409 {object} response.Response "Conflict - already active or profile incomplete"
// @Router /stores/my/activate [put]
func (h *StoreHandler) ActivateStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	err := h.storeUsecase.ActivateStore(userID)
	if err != nil {
		if err.Error() == "STORE_SUSPENDED_BY_ADMIN" {
			return response.Forbidden(c, err.Error())
		}
		if err.Error() == "STORE_ALREADY_ACTIVE" || err.Error() == "STORE_PROFILE_INCOMPLETE" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store activated successfully", nil)
}

// DeactivateStore godoc
// @Summary Deactivate store (Seller only)
// @Description Deactivate the authenticated user's store. Only active stores can be deactivated.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Store deactivated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 409 {object} response.Response "Conflict - store not active"
// @Router /stores/my/deactivate [put]
func (h *StoreHandler) DeactivateStore(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	err := h.storeUsecase.DeactivateStore(userID)
	if err != nil {
		if err.Error() == "STORE_NOT_ACTIVE" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store deactivated successfully", nil)
}

// SuspendStore godoc
// @Summary Suspend a store (Admin only)
// @Description Suspend any store. Admin can suspend stores for violations, reports, etc.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Store ID"
// @Success 200 {object} response.Response "Store suspended successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - already suspended"
// @Router /admin/stores/{id}/suspend [put]
func (h *StoreHandler) SuspendStore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid store ID")
	}

	err = h.storeUsecase.SuspendStore(id)
	if err != nil {
		if err.Error() == "STORE_ALREADY_SUSPENDED" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store suspended successfully", nil)
}

// UnsuspendStore godoc
// @Summary Unsuspend a store (Admin only)
// @Description Unsuspend a suspended store. Store will be set to inactive, seller must activate it.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Store ID"
// @Success 200 {object} response.Response "Store unsuspended successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - not suspended"
// @Router /admin/stores/{id}/unsuspend [put]
func (h *StoreHandler) UnsuspendStore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid store ID")
	}

	err = h.storeUsecase.UnsuspendStore(id)
	if err != nil {
		if err.Error() == "STORE_NOT_SUSPENDED" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store unsuspended successfully", nil)
}
// ApproveStore godoc
// @Summary Approve pending store (Admin only)
// @Description Approve a pending store to make it active. Admin can approve stores waiting for verification.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Store ID"
// @Success 200 {object} response.Response "Store approved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - store not pending"
// @Router /admin/stores/{id}/approve [put]
func (h *StoreHandler) ApproveStore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid store ID")
	}

	err = h.storeUsecase.ApproveStore(id)
	if err != nil {
		if err.Error() == "STORE_NOT_PENDING" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store approved successfully", nil)
}

// RejectStore godoc
// @Summary Reject pending store (Admin only)
// @Description Reject a pending store. Store will be set to inactive status.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Store ID"
// @Success 200 {object} response.Response "Store rejected successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 409 {object} response.Response "Conflict - store not pending"
// @Router /admin/stores/{id}/reject [put]
func (h *StoreHandler) RejectStore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid store ID")
	}

	err = h.storeUsecase.RejectStore(id)
	if err != nil {
		if err.Error() == "STORE_NOT_PENDING" {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Store rejected successfully", nil)
}

// GetPendingStores godoc
// @Summary Get pending stores (Admin only)
// @Description Get all stores waiting for admin approval with pagination and search.
// @Tags Stores
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by store name"
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Store} "Pending stores retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /admin/stores/pending [get]
func (h *StoreHandler) GetPendingStores(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	stores, meta, err := h.storeUsecase.GetPendingStores(page, limit, search)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Pending stores retrieved successfully", stores, meta)
}