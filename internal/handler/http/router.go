package http

import (
	"go-commerce/internal/domain"
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
	users.Get("/profile", middleware.JWTMiddleware(r.jwtManager), userHandler.GetProfile)
	users.Put("/my", middleware.JWTMiddleware(r.jwtManager), userHandler.UpdateProfile)
	users.Put("/profile", middleware.JWTMiddleware(r.jwtManager), userHandler.UpdateProfile)
	users.Put("/my/password", middleware.JWTMiddleware(r.jwtManager), userHandler.ChangePassword)
}

func (r *Router) SetupStoreRoutes(storeUsecase *usecase.StoreUsecase) {
	storeHandler := NewStoreHandler(storeUsecase)
	
	api := r.app.Group("/api/v1")
	stores := api.Group("/stores")

	// Public routes - only show active stores
	stores.Get("/", storeHandler.GetAllStores)

	// Protected routes
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	stores.Post("/", jwtMiddleware, storeHandler.CreateStore)
	stores.Get("/my", jwtMiddleware, storeHandler.GetMyStore)
	stores.Put("/my", jwtMiddleware, storeHandler.UpdateMyStore)

	// Store status management (seller only)
	stores.Put("/my/activate", jwtMiddleware, storeHandler.ActivateStore)
	stores.Put("/my/deactivate", jwtMiddleware, storeHandler.DeactivateStore)

	// Public route with ID parameter (must be after /my routes)
	stores.Get("/:id", storeHandler.GetStoreByID)

	// Admin routes
	admin := api.Group("/admin")
	adminMiddleware := middleware.JWTMiddleware(r.jwtManager)
	requireAdmin := middleware.RequireAdmin()
	admin.Get("/stores/pending", adminMiddleware, requireAdmin, storeHandler.GetPendingStores)
	admin.Put("/stores/:id/approve", adminMiddleware, requireAdmin, storeHandler.ApproveStore)
	admin.Put("/stores/:id/reject", adminMiddleware, requireAdmin, storeHandler.RejectStore)
	admin.Put("/stores/:id/suspend", adminMiddleware, requireAdmin, storeHandler.SuspendStore)
	admin.Put("/stores/:id/unsuspend", adminMiddleware, requireAdmin, storeHandler.UnsuspendStore)
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
	categories.Put("/:id/activate", adminMiddleware, requireAdmin, categoryHandler.ActivateCategory)
	categories.Put("/:id/deactivate", adminMiddleware, requireAdmin, categoryHandler.DeactivateCategory)
	categories.Delete("/:id", adminMiddleware, requireAdmin, categoryHandler.DeleteCategory)
}

func (r *Router) SetupAddressRoutes(addressUsecase *usecase.AddressUsecase) {
	addressHandler := NewAddressHandler(addressUsecase)
	
	api := r.app.Group("/api/v1")
	addresses := api.Group("/addresses")

	// Protected routes (user can only manage their own addresses)
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	addresses.Get("/", jwtMiddleware, addressHandler.GetMyAddresses)
	addresses.Get("/default", jwtMiddleware, addressHandler.GetDefaultAddress)
	addresses.Post("/", jwtMiddleware, addressHandler.CreateAddress)
	addresses.Get("/:id", jwtMiddleware, addressHandler.GetAddressByID)
	addresses.Put("/:id", jwtMiddleware, addressHandler.UpdateAddress)
	addresses.Put("/:id/default", jwtMiddleware, addressHandler.SetDefaultAddress)
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
	products.Get("/status", middleware.JWTMiddleware(r.jwtManager), middleware.RequireAdmin(), productHandler.GetProductsByStatus)
	products.Get("/search/slug", productHandler.SearchProductsBySlug)
	products.Get("/slug/:slug", productHandler.GetProductBySlug)

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

	// Product status management (seller only)
	products.Put("/:id/activate", jwtMiddleware, productHandler.ActivateProduct)
	products.Put("/:id/deactivate", jwtMiddleware, productHandler.DeactivateProduct)

	// Public route with ID parameter (must be after /my routes)
	products.Get("/:id", productHandler.GetProductByID)

	// Admin routes
	admin := api.Group("/admin")
	adminMiddleware := middleware.JWTMiddleware(r.jwtManager)
	requireAdmin := middleware.RequireAdmin()
	admin.Put("/products/:id/suspend", adminMiddleware, requireAdmin, productHandler.SuspendProduct)
	admin.Put("/products/:id/unsuspend", adminMiddleware, requireAdmin, productHandler.UnsuspendProduct)
}

func (r *Router) SetupTransactionRoutes(transactionUsecase *usecase.TransactionUsecase, paymentIntentUsecase domain.PaymentIntentUsecase) {
	transactionHandler := NewTransactionHandler(transactionUsecase, paymentIntentUsecase)
	
	api := r.app.Group("/api/v1")
	transactions := api.Group("/transactions")

	// Protected routes - Buyer operations
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	transactions.Post("/", jwtMiddleware, transactionHandler.CreateTransaction)
	transactions.Get("/my", jwtMiddleware, transactionHandler.GetMyTransactions)
	transactions.Put("/:id/confirm-delivery", jwtMiddleware, transactionHandler.ConfirmDelivered)
	transactions.Put("/:id/cancel", jwtMiddleware, transactionHandler.CancelTransaction)

	// Transaction by ID (must be after /my routes)
	transactions.Get("/:id", jwtMiddleware, transactionHandler.GetTransactionByID)

	// Seller operations
	seller := api.Group("/seller")
	seller.Put("/transactions/:id/process", jwtMiddleware, transactionHandler.ProcessOrder)
	seller.Put("/transactions/:id/ship", jwtMiddleware, transactionHandler.ShipOrder)

	// Admin operations
	admin := api.Group("/admin")
	adminMiddleware := middleware.JWTMiddleware(r.jwtManager)
	requireAdmin := middleware.RequireAdmin()
	admin.Put("/transactions/:id/refund", adminMiddleware, requireAdmin, transactionHandler.RefundTransaction)
}

func (r *Router) SetupPaymentIntentRoutes(paymentIntentUsecase domain.PaymentIntentUsecase) {
	paymentIntentHandler := NewPaymentIntentHandler(paymentIntentUsecase)
	
	api := r.app.Group("/api/v1")
	transactions := api.Group("/transactions")

	// Payment intent creation (buyer)
	jwtMiddleware := middleware.JWTMiddleware(r.jwtManager)
	transactions.Post("/:id/pay", jwtMiddleware, paymentIntentHandler.CreatePaymentIntent)

	// Admin payment simulation endpoints
	admin := api.Group("/admin")
	adminMiddleware := middleware.JWTMiddleware(r.jwtManager)
	requireAdmin := middleware.RequireAdmin()
	admin.Put("/payments/:intentId/simulate-success", adminMiddleware, requireAdmin, paymentIntentHandler.SimulatePaymentSuccess)
	admin.Put("/payments/:intentId/simulate-failed", adminMiddleware, requireAdmin, paymentIntentHandler.SimulatePaymentFailed)

	// Payment gateway callback (updates payment intent)
	callbacks := api.Group("/callbacks")
	callbacks.Post("/payments/:intentId", paymentIntentHandler.OnPaymentCallback)
}