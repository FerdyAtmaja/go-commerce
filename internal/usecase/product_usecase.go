package usecase

import (
	"errors"

	"go-commerce/internal/domain"
	"go-commerce/pkg/utils"
)

type ProductUsecase struct {
	productRepo  domain.ProductRepository
	photoRepo    domain.PhotoProdukRepository
	storeRepo    domain.StoreRepository
	categoryRepo domain.CategoryRepository
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	photoRepo domain.PhotoProdukRepository,
	storeRepo domain.StoreRepository,
	categoryRepo domain.CategoryRepository,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:  productRepo,
		photoRepo:    photoRepo,
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *ProductUsecase) CreateProduct(userID uint64, req *domain.CreateProductRequest) (*domain.Product, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("store not found")
	}

	// Validate category exists and is active
	category, err := u.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if category.Status != "active" {
		return nil, errors.New("invalid category")
	}
	if !category.IsLeaf {
		return nil, errors.New("product must use leaf category")
	}

	// Generate unique slug from product name
	baseSlug := utils.GenerateSlug(req.NamaProduk)
	slug := utils.EnsureUniqueSlug(baseSlug, func(s string) bool {
		_, err := u.productRepo.GetBySlug(s)
		return err == nil // true if slug exists
	})

	product := &domain.Product{
		NamaProduk:    req.NamaProduk,
		Slug:          slug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IDToko:        store.ID,
		IDCategory:    req.IDCategory,
		Status:        getProductStatus(req.Status),
		Berat:         req.Berat,
		SoldCount:     0,
	}

	err = u.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	// Update category's has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(req.IDCategory)
	}()

	// Get the created product with all relations
	createdProduct, err := u.productRepo.GetByID(product.ID)
	if err != nil {
		return nil, err
	}

	// Update search index async
	go func() {
		// Simulate updating search index
		// updateSearchIndex(product.ID, product.NamaProduk)
	}()

	return createdProduct, nil
}

func (u *ProductUsecase) GetProductByID(id uint64) (*domain.Product, error) {
	// Track product view async
	go func() {
		// Simulate tracking product view analytics
		// trackProductView(id)
	}()

	return u.productRepo.GetByID(id)
}

func (u *ProductUsecase) GetProductBySlug(slug string) (*domain.Product, error) {
	return u.productRepo.GetBySlug(slug)
}

func (u *ProductUsecase) SearchProductsBySlug(slugPattern string, page, limit int) ([]*domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.productRepo.SearchBySlug(slugPattern, limit, offset)
}

func (u *ProductUsecase) GetMyProducts(userID uint64, page, limit int, search string) ([]*domain.Product, int64, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return nil, 0, errors.New("store not found")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.productRepo.GetByTokoID(store.ID, limit, offset, search)
}

func (u *ProductUsecase) GetAllProducts(filter *domain.ProductFilter) ([]*domain.Product, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}

	// Use advanced filtering if available, otherwise fallback to basic
	if filter.MinPrice != "" || filter.MaxPrice != "" ||
		(filter.SortBy != "" && filter.SortBy != "newest") {
		return u.productRepo.GetAllWithFilter(filter)
	}

	// Fallback to basic filtering
	offset := (filter.Page - 1) * filter.Limit
	return u.productRepo.GetAll(filter.Limit, offset, filter.Search, filter.CategoryID)
}

func (u *ProductUsecase) UpdateProduct(userID, productID uint64, req *domain.UpdateProductRequest) (*domain.Product, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return nil, err
	}

	// Get existing product (use management method for seller access)
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return nil, err
	}

	// Validate category exists and is active
	category, err := u.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if category.Status != "active" {
		return nil, errors.New("invalid category")
	}
	if !category.IsLeaf {
		return nil, errors.New("product must use leaf category")
	}

	// Check if name changed to update slug
	nameChanged := product.NamaProduk != req.NamaProduk
	categoryChanged := product.IDCategory != req.IDCategory
	oldCategoryID := product.IDCategory

	// Update fields
	product.NamaProduk = req.NamaProduk
	product.HargaReseller = req.HargaReseller
	product.HargaKonsumen = req.HargaKonsumen
	product.Stok = req.Stok
	product.Deskripsi = req.Deskripsi
	product.IDCategory = req.IDCategory
	product.Berat = req.Berat
	if req.Status != "" {
		product.Status = req.Status
	}

	// Update slug if name changed
	if nameChanged {
		baseSlug := utils.GenerateSlug(req.NamaProduk)
		slug := utils.EnsureUniqueSlug(baseSlug, func(s string) bool {
			existingProduct, err := u.productRepo.GetBySlug(s)
			return err == nil && existingProduct.ID != productID // true if slug exists and not current product
		})
		product.Slug = slug
	}

	err = u.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	// Update category flags if category changed
	if categoryChanged {
		go func() {
			u.categoryRepo.UpdateHasActiveProduct(oldCategoryID)
			u.categoryRepo.UpdateHasActiveProduct(req.IDCategory)
		}()
	}

	// Return updated product with all relations (use management method)
	return u.productRepo.GetByIDForManagement(productID)
}

func (u *ProductUsecase) DeleteProduct(userID, productID uint64) error {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return err
	}

	// Get product to know its category before deletion
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return err
	}
	categoryID := product.IDCategory

	err = u.productRepo.Delete(productID)
	if err != nil {
		return err
	}

	// Update category's has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(categoryID)
	}()

	return nil
}

func (u *ProductUsecase) AddProductPhoto(userID, productID uint64, photoURL string, isPrimary bool) (*domain.PhotoProduk, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return nil, err
	}

	// Get existing photos to determine position
	existingPhotos, err := u.photoRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	position := int64(len(existingPhotos) + 1)

	photo := &domain.PhotoProduk{
		IDProduk:  productID,
		URL:       photoURL,
		IsPrimary: isPrimary,
		Position:  position,
	}

	err = u.photoRepo.Create(photo)
	if err != nil {
		return nil, err
	}

	// If this is set as primary, update other photos
	if isPrimary {
		err = u.photoRepo.SetPrimary(productID, photo.ID)
		if err != nil {
			return nil, err
		}
	}

	return photo, nil
}

func (u *ProductUsecase) SetPrimaryPhoto(userID, productID, photoID uint64) error {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return err
	}

	return u.photoRepo.SetPrimary(productID, photoID)
}

func (u *ProductUsecase) DeleteProductPhoto(userID, productID, photoID uint64) error {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return err
	}

	return u.photoRepo.Delete(photoID)
}

// Helper function to get product status with default
func getProductStatus(status string) string {
	if status == "" {
		return "active"
	}
	return status
}

// GetProductsByStatus gets products by status (admin only)
func (u *ProductUsecase) GetProductsByStatus(status string, page, limit int) ([]*domain.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.productRepo.GetByStatus(status, limit, offset)
}

// ActivateProduct implements business rules for product activation
func (u *ProductUsecase) ActivateProduct(userID, productID uint64) error {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return err
	}

	// Get current product
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Business rules from rancangan-update.txt
	if product.Status == "active" {
		return errors.New("product already active") // 409 noop
	}

	// Check store status
	if store.Status != "active" {
		return errors.New("cannot activate product: store inactive") // 403 forbidden
	}

	// Check category status
	category, err := u.categoryRepo.GetByID(product.IDCategory)
	if err != nil {
		return errors.New("category not found")
	}
	if category.Status != "active" {
		return errors.New("cannot activate product: category inactive") // 409 conflict
	}

	// Check stock (optional reject based on business requirement)
	if product.Stok <= 0 {
		return errors.New("cannot activate product: out of stock") // 409 conflict
	}

	// Update status
	product.Status = "active"
	if err := u.productRepo.Update(product); err != nil {
		return errors.New("failed to activate product")
	}

	// Update category has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(product.IDCategory)
	}()

	return nil
}

// DeactivateProduct implements business rules for product deactivation
func (u *ProductUsecase) DeactivateProduct(userID, productID uint64) error {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	// Check ownership
	err = u.productRepo.CheckOwnership(productID, store.ID)
	if err != nil {
		return err
	}

	// Get current product
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Business rules from rancangan-update.txt
	if product.Status == "inactive" {
		return errors.New("product already inactive") // 409 noop
	}

	// Deactivation is always allowed for sellers
	product.Status = "inactive"
	if err := u.productRepo.Update(product); err != nil {
		return errors.New("failed to deactivate product")
	}

	// Update category has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(product.IDCategory)
	}()

	return nil
}

// SuspendProduct allows admin to force deactivate any product
func (u *ProductUsecase) SuspendProduct(productID uint64) error {
	// Get product
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Status == "suspended" {
		return errors.New("product already suspended")
	}

	// Admin can force suspend
	product.Status = "suspended"
	if err := u.productRepo.Update(product); err != nil {
		return errors.New("failed to suspend product")
	}

	// Update category has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(product.IDCategory)
	}()

	return nil
}

// UnsuspendProduct allows admin to reactivate suspended products
func (u *ProductUsecase) UnsuspendProduct(productID uint64) error {
	// Get product
	product, err := u.productRepo.GetByIDForManagement(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Status != "suspended" {
		return errors.New("product is not suspended")
	}

	// Reactivate to active status
	product.Status = "active"
	if err := u.productRepo.Update(product); err != nil {
		return errors.New("failed to unsuspend product")
	}

	// Update category has_active_product flag
	go func() {
		u.categoryRepo.UpdateHasActiveProduct(product.IDCategory)
	}()

	return nil
}
