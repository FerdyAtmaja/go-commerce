package http

import (
	"go-commerce/internal/handler/middleware"
	"go-commerce/internal/usecase"
	"go-commerce/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app        *fiber.App
	jwtManager *jwt.JWTManager
}

func NewRouter(app *fiber.App, jwtManager *jwt.JWTManager) *Router {
	return &Router{
		app:        app,
		jwtManager: jwtManager,
	}
}

func (r *Router) SetupAuthRoutes(authUsecase *usecase.AuthUsecase) {
	authHandler := NewAuthHandler(authUsecase)
	
	api := r.app.Group("/api/v1")
	auth := api.Group("/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	auth.Post("/logout", middleware.JWTMiddleware(r.jwtManager), authHandler.Logout)
	
	// User profile routes
	users := api.Group("/users")
	users.Get("/my", middleware.JWTMiddleware(r.jwtManager), authHandler.GetProfile)
}