package http

import (
	"go-commerce/internal/domain"
	"go-commerce/internal/handler/middleware"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	validator   *validator.Validate
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		validator:   validator.New(),
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Profile retrieved successfully", user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	user, err := h.userUsecase.UpdateProfile(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Profile updated successfully", user)
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req domain.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	if err := h.userUsecase.ChangePassword(userID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Password changed successfully", nil)
}