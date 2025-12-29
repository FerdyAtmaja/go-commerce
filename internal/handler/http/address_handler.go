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

type AddressHandler struct {
	addressUsecase *usecase.AddressUsecase
	validator      *validator.Validate
}

func NewAddressHandler(addressUsecase *usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		addressUsecase: addressUsecase,
		validator:      validator.New(),
	}
}

// CreateAddress godoc
// @Summary Create a new address (Authenticated User)
// @Description Create a new address for the authenticated user. Requires authentication.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateAddressRequest true "Address creation request"
// @Success 201 {object} response.Response{data=domain.Address} "Address created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /addresses [post]
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req domain.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	address, err := h.addressUsecase.CreateAddress(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Address created successfully", address)
}

// GetMyAddresses godoc
// @Summary Get current user's addresses (Authenticated User)
// @Description Get all addresses for the authenticated user with pagination. Requires authentication.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Address} "Addresses retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /addresses [get]
func (h *AddressHandler) GetMyAddresses(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	addresses, meta, err := h.addressUsecase.GetMyAddresses(userID, page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Addresses retrieved successfully", addresses, meta)
}

// GetAddressByID godoc
// @Summary Get address by ID (Authenticated User)
// @Description Get a single address by its ID (owner only). Only address owner can access.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} response.Response{data=domain.Address} "Address retrieved successfully"
// @Failure 400 {object} response.Response "Invalid address ID"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Address not found"
// @Router /addresses/{id} [get]
func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	address, err := h.addressUsecase.GetAddressByID(addressID, userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Address retrieved successfully", address)
}

// UpdateAddress godoc
// @Summary Update an address (Authenticated User)
// @Description Update an existing address (owner only). Only address owner can update.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param request body domain.UpdateAddressRequest true "Address update request"
// @Success 200 {object} response.Response{data=domain.Address} "Address updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Address not found"
// @Router /addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	var req domain.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	address, err := h.addressUsecase.UpdateAddress(addressID, userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Address updated successfully", address)
}

// DeleteAddress godoc
// @Summary Delete an address (Authenticated User)
// @Description Delete an existing address (owner only). Only address owner can delete.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} response.Response "Address deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Address not found"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	if err := h.addressUsecase.DeleteAddress(addressID, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Address deleted successfully", nil)
}

// SetDefaultAddress godoc
// @Summary Set default address (Authenticated User)
// @Description Set an address as default for the authenticated user. Requires authentication.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} response.Response "Default address set successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Address not found"
// @Router /addresses/{id}/default [put]
func (h *AddressHandler) SetDefaultAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	if err := h.addressUsecase.SetDefaultAddress(addressID, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Default address set successfully", nil)
}

// GetDefaultAddress godoc
// @Summary Get default address (Authenticated User)
// @Description Get the default address for the authenticated user. Requires authentication.
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.Address} "Default address retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "No default address found"
// @Router /addresses/default [get]
func (h *AddressHandler) GetDefaultAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	address, err := h.addressUsecase.GetDefaultAddress(userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Default address retrieved successfully", address)
}

// GetProvinces godoc
// @Summary Get all provinces (Public)
// @Description Get all provinces in Indonesia. This is a public endpoint accessible to everyone.
// @Tags Regions
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]domain.Province} "Provinces retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /regions/provinces [get]
func (h *AddressHandler) GetProvinces(c *fiber.Ctx) error {
	provinces, err := h.addressUsecase.GetProvinces()
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "Provinces retrieved successfully", provinces)
}

// GetCitiesByProvince godoc
// @Summary Get cities by province (Public)
// @Description Get all cities in a specific province. This is a public endpoint accessible to everyone.
// @Tags Regions
// @Accept json
// @Produce json
// @Param provinceId path string true "Province ID"
// @Success 200 {object} response.Response{data=[]domain.City} "Cities retrieved successfully"
// @Failure 400 {object} response.Response "Province ID is required"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /regions/provinces/{provinceId}/cities [get]
func (h *AddressHandler) GetCitiesByProvince(c *fiber.Ctx) error {
	provinceID := c.Params("provinceId")
	if provinceID == "" {
		return response.BadRequest(c, "Province ID is required")
	}

	cities, err := h.addressUsecase.GetCitiesByProvince(provinceID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "Cities retrieved successfully", cities)
}
