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

func (r *Router) SetupCategoryRoutes(categoryUsecase *usecase.CategoryUsecase) {
	categoryHandler := NewCategoryHandler(categoryUsecase)
	
	api := r.app.Group("/api/v1")
	categories := api.Group("/categories")

	// Public routes
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/:id", categoryHandler.GetCategoryByID)

	// Admin only routes
	adminMiddleware := middleware.JWTMiddleware(r.jwtManager)
	requireAdmin := middleware.RequireAdmin()
	categories.Post("/", adminMiddleware, requireAdmin, categoryHandler.CreateCategory)
	categories.Put("/:id", adminMiddleware, requireAdmin, categoryHandler.UpdateCategory)
	categories.Delete("/:id", adminMiddleware, requireAdmin, categoryHandler.DeleteCategory)
}

func (r *Router) SetupAddressRoutes(addressUsecase *usecase.AddressUsecase) {
	addressHandler := NewAddressHandler(addressUsecase)
	
	api := r.app.Group("/api/v1")
	addresses := api.Group("/addresses")

	// Protected routes (user can only manage their own addresses)
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	addresses.Get("/", jwtMiddleware, addressHandler.GetMyAddresses)
	addresses.Post("/", jwtMiddleware, addressHandler.CreateAddress)
	addresses.Get("/:id", jwtMiddleware, addressHandler.GetAddressByID)
	addresses.Put("/:id", jwtMiddleware, addressHandler.UpdateAddress)
	addresses.Delete("/:id", jwtMiddleware, addressHandler.DeleteAddress)

	// Utility routes for Indonesia regions
	provinces := api.Group("/provinces")
	provinces.Get("/", addressHandler.GetProvinces)
	provinces.Get("/:provinceId/cities", addressHandler.GetCitiesByProvince)
}

func (r *Router) SetupProductRoutes(productUsecase *usecase.ProductUsecase) {
	productHandler := NewProductHandler(productUsecase)
	
	api := r.app.Group("/api/v1")
	products := api.Group("/products")

	// Public routes
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/slug/:slug", productHandler.GetProductBySlug)
	products.Get("/:id", productHandler.GetProductByID)

	// Protected routes (store owner only)
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	products.Get("/my", jwtMiddleware, productHandler.GetMyProducts)
	products.Post("/", jwtMiddleware, productHandler.CreateProduct)
	products.Put("/:id", jwtMiddleware, productHandler.UpdateProduct)
	products.Delete("/:id", jwtMiddleware, productHandler.DeleteProduct)

	// Photo management routes
	products.Post("/:id/photos", jwtMiddleware, productHandler.UploadProductPhoto)
	products.Put("/:id/photos/:photoId/primary", jwtMiddleware, productHandler.SetPrimaryPhoto)
	products.Delete("/:id/photos/:photoId", jwtMiddleware, productHandler.DeleteProductPhoto)
}

func (r *Router) SetupTransactionRoutes(transactionUsecase *usecase.TransactionUsecase) {
	transactionHandler := NewTransactionHandler(transactionUsecase)
	
	api := r.app.Group("/api/v1")
	transactions := api.Group("/transactions")

	// Protected routes
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	transactions.Post("/", jwtMiddleware, transactionHandler.CreateTransaction)
	transactions.Get("/my", jwtMiddleware, transactionHandler.GetMyTransactions)
	transactions.Get("/:id", jwtMiddleware, transactionHandler.GetTransactionByID)
}