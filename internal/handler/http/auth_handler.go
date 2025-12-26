package http

import (
	"go-commerce/internal/domain"
	"go-commerce/internal/handler/middleware"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	validator   *validator.Validate
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	authResponse, err := h.authUsecase.Register(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "User registered successfully", authResponse)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed: "+err.Error())
	}

	authResponse, err := h.authUsecase.Login(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Login successful", authResponse)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	user, err := h.authUsecase.GetUserByID(userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Profile retrieved successfully", user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	return response.Success(c, "Logout successful", nil)
}