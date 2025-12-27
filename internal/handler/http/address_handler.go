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

func (h *AddressHandler) GetAddressByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	address, err := h.addressUsecase.GetAddressByID(uint(addressID), userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Address retrieved successfully", address)
}

func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	var req domain.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	address, err := h.addressUsecase.UpdateAddress(uint(addressID), userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Address updated successfully", address)
}

func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	idParam := c.Params("id")
	addressID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid address ID")
	}

	if err := h.addressUsecase.DeleteAddress(uint(addressID), userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Address deleted successfully", nil)
}

// Indonesia region endpoints
func (h *AddressHandler) GetProvinces(c *fiber.Ctx) error {
	provinces, err := h.addressUsecase.GetProvinces()
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "Provinces retrieved successfully", provinces)
}

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