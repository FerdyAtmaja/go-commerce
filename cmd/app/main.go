package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/http"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/repository/mysql"
	"go-commerce/internal/usecase"
	"go-commerce/pkg/config"
	"go-commerce/pkg/database"
	"go-commerce/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg.JWT.RefreshExpireHours)

	// Initialize repositories
	userRepo := mysql.NewUserRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtManager)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(response.Response{
				Status:  "error",
				Message: err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, "Server is running", fiber.Map{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Setup routes
	router := http.NewRouter(app, jwtManager)
	router.SetupAuthRoutes(authUsecase)

	// API info endpoint
	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, "Go Commerce API v1", fiber.Map{
			"version": "1.0.0",
			"message": "Welcome to Go Commerce API",
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// Close database connection
	sqlDB, _ := db.DB()
	sqlDB.Close()
	log.Println("Server stopped")
}