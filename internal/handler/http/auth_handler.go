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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, and profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration request"
// @Success 201 {object} response.Response{data=domain.AuthResponse} "User registered successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=domain.AuthResponse} "Login successful"
// @Failure 400 {object} response.Response "Bad request"
// @Router /auth/login [post]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.User} "Profile retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "User not found"
// @Router /users/my [get]
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

// Logout godoc
// @Summary User logout
// @Description Logout current user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Logout successful"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	return response.Success(c, "Logout successful", nil)
}