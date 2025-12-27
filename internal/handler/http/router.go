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
}

func (r *Router) SetupUserRoutes(userUsecase *usecase.UserUsecase) {
	userHandler := NewUserHandler(userUsecase)
	
	api := r.app.Group("/api/v1")
	users := api.Group("/users")

	// Protected routes
	users.Get("/my", middleware.JWTMiddleware(r.jwtManager), userHandler.GetProfile)
	users.Put("/my", middleware.JWTMiddleware(r.jwtManager), userHandler.UpdateProfile)
	users.Put("/my/password", middleware.JWTMiddleware(r.jwtManager), userHandler.ChangePassword)
}

func (r *Router) SetupStoreRoutes(storeUsecase *usecase.StoreUsecase) {
	storeHandler := NewStoreHandler(storeUsecase)
	
	api := r.app.Group("/api/v1")
	stores := api.Group("/stores")

	// Public routes
	stores.Get("/", storeHandler.GetAllStores)
	stores.Get("/:id", storeHandler.GetStoreByID)

	// Protected routes
	stores.Get("/my", middleware.JWTMiddleware(r.jwtManager), storeHandler.GetMyStore)
	stores.Put("/my", middleware.JWTMiddleware(r.jwtManager), storeHandler.UpdateMyStore)
}