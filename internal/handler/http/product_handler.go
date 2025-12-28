package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/middleware"
	"go-commerce/internal/handler/response"
	"go-commerce/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	validator      *validator.Validate
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		validator:      validator.New(),
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for the authenticated user's store
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateProductRequest true "Product creation request"
// @Success 201 {object} response.Response{data=domain.Product} "Product created successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req domain.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed")
	}

	userID := middleware.GetUserID(c)
	product, err := h.productUsecase.CreateProduct(userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Product created successfully", product)
}

// GetMyProducts godoc
// @Summary Get current user's products
// @Description Get all products owned by the authenticated user with pagination and search
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by product name"
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Product} "Products retrieved successfully"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /products/my [get]
func (h *ProductHandler) GetMyProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	userID := middleware.GetUserID(c)
	products, total, err := h.productUsecase.GetMyProducts(userID, page, limit, search)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Products retrieved successfully", products, response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: (int(total) + limit - 1) / limit,
	})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get all products with pagination and filtering (public endpoint)
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by product name"
// @Param category_id query string false "Filter by category ID"
// @Param min_price query string false "Minimum price filter"
// @Param max_price query string false "Maximum price filter"
// @Param sort_by query string false "Sort by: newest, oldest, price_asc, price_desc, popular, name_asc, name_desc" default(newest)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Product} "Products retrieved successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	filter := &domain.ProductFilter{
		Search:     c.Query("search", ""),
		CategoryID: c.Query("category_id", ""),
		MinPrice:   c.Query("min_price", ""),
		MaxPrice:   c.Query("max_price", ""),
		SortBy:     c.Query("sort_by", "newest"),
		Page:       page,
		Limit:      limit,
	}

	products, total, err := h.productUsecase.GetAllProducts(filter)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Products retrieved successfully", products, response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: (int(total) + limit - 1) / limit,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get a single product by its ID (public endpoint)
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response{data=domain.Product} "Product retrieved successfully"
// @Failure 400 {object} response.Response "Invalid product ID"
// @Failure 404 {object} response.Response "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	product, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Product retrieved successfully", product)
}

// GetProductBySlug godoc
// @Summary Get product by slug
// @Description Get a single product by its exact slug (public endpoint)
// @Tags Products
// @Accept json
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} response.Response{data=domain.Product} "Product retrieved successfully"
// @Failure 400 {object} response.Response "Invalid slug"
// @Failure 404 {object} response.Response "Product not found"
// @Router /products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.BadRequest(c, "Slug is required")
	}

	product, err := h.productUsecase.GetProductBySlug(slug)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.Success(c, "Product retrieved successfully", product)
}

// SearchProductsBySlug godoc
// @Summary Search products by slug pattern
// @Description Search products by partial slug match with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param slug query string true "Slug pattern to search"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Product} "Products found successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Router /products/search/slug [get]
func (h *ProductHandler) SearchProductsBySlug(c *fiber.Ctx) error {
	slugPattern := c.Query("slug")
	if slugPattern == "" {
		return response.BadRequest(c, "Slug parameter is required")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, total, err := h.productUsecase.SearchProductsBySlug(slugPattern, page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Products found successfully", products, response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: (int(total) + limit - 1) / limit,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body domain.UpdateProductRequest true "Product update request"
// @Success 200 {object} response.Response{data=domain.Product} "Product updated successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - not product owner"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	var req domain.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "Validation failed")
	}

	userID := middleware.GetUserID(c)
	product, err := h.productUsecase.UpdateProduct(userID, id, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete an existing product (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response "Product deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - not product owner"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	userID := middleware.GetUserID(c)
	err = h.productUsecase.DeleteProduct(userID, id)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Product deleted successfully", nil)
}

// UploadProductPhoto godoc
// @Summary Upload product photo
// @Description Upload a photo for a product (owner only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param photo formData file true "Photo file (JPG, JPEG, PNG, max 5MB)"
// @Param is_primary formData boolean false "Set as primary photo" default(false)
// @Success 201 {object} response.Response{data=domain.PhotoProduk} "Photo uploaded successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - not product owner"
// @Router /products/{id}/photos [post]
func (h *ProductHandler) UploadProductPhoto(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	// Get file from form
	file, err := c.FormFile("photo")
	if err != nil {
		return response.BadRequest(c, "Photo file is required")
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return response.BadRequest(c, "Invalid file type. Only JPG, JPEG, PNG allowed")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return response.BadRequest(c, "File size too large. Maximum 5MB allowed")
	}

	// Generate unique filename
	filename := generateFileName(file.Filename)
	
	// Create upload directory if not exists
	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return response.InternalServerError(c, "Failed to create upload directory")
	}

	// Save file
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return response.InternalServerError(c, "Failed to save file")
	}

	// Process image async (resize, thumbnail)
	go func() {
		// Simulate image processing
		// resizeImage(filePath)
		// generateThumbnail(filePath)
	}()

	// Get isPrimary from form
	isPrimary := c.FormValue("is_primary") == "true"

	// Add photo to database
	userID := middleware.GetUserID(c)
	photoURL := fmt.Sprintf("/uploads/products/%s", filename)
	
	photo, err := h.productUsecase.AddProductPhoto(userID, productID, photoURL, isPrimary)
	if err != nil {
		// Delete uploaded file if database operation fails
		os.Remove(filePath)
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Photo uploaded successfully", photo)
}

// SetPrimaryPhoto godoc
// @Summary Set primary photo
// @Description Set a photo as the primary photo for a product (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param photoId path int true "Photo ID"
// @Success 200 {object} response.Response "Primary photo set successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - not product owner"
// @Router /products/{id}/photos/{photoId}/primary [put]
func (h *ProductHandler) SetPrimaryPhoto(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	photoID, err := strconv.ParseUint(c.Params("photoId"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid photo ID")
	}

	userID := middleware.GetUserID(c)
	err = h.productUsecase.SetPrimaryPhoto(userID, productID, photoID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Primary photo set successfully", nil)
}

// DeleteProductPhoto godoc
// @Summary Delete product photo
// @Description Delete a photo from a product (owner only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param photoId path int true "Photo ID"
// @Success 200 {object} response.Response "Photo deleted successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - not product owner"
// @Router /products/{id}/photos/{photoId} [delete]
func (h *ProductHandler) DeleteProductPhoto(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	photoID, err := strconv.ParseUint(c.Params("photoId"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid photo ID")
	}

	userID := middleware.GetUserID(c)
	err = h.productUsecase.DeleteProductPhoto(userID, productID, photoID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Photo deleted successfully", nil)
}

// GetProductsByStatus godoc
// @Summary Get products by status (Admin only)
// @Description Get products filtered by status with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string true "Product status: active or inactive"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]domain.Product} "Products retrieved successfully"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden - admin only"
// @Router /products/status [get]
func (h *ProductHandler) GetProductsByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	if status == "" {
		return response.BadRequest(c, "Status parameter is required")
	}

	if status != "active" && status != "inactive" {
		return response.BadRequest(c, "Status must be 'active' or 'inactive'")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, total, err := h.productUsecase.GetProductsByStatus(status, page, limit)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Paginated(c, "Products retrieved successfully", products, response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: (int(total) + limit - 1) / limit,
	})
}

// Helper functions
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg", 
		"image/png",
	}
	
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)
	
	// Generate UUID for uniqueness
	uuid := uuid.New().String()
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s_%d_%s%s", name, timestamp, uuid[:8], ext)
}