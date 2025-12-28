package usecase

import (
	"errors"
	"fmt"
	"strings"

	"go-commerce/internal/domain"
)

type ProductUsecase struct {
	productRepo      domain.ProductRepository
	photoRepo        domain.PhotoProdukRepository
	storeRepo        domain.StoreRepository
	categoryRepo     domain.CategoryRepository
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

	// Validate category exists
	_, err = u.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Generate slug from product name
	slug := generateSlug(req.NamaProduk)
	
	// Check if slug already exists, if yes, append number
	originalSlug := slug
	counter := 1
	for {
		_, err := u.productRepo.GetBySlug(slug)
		if err != nil {
			// Slug not found, we can use it
			break
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	product := &domain.Product{
		NamaProduk:    req.NamaProduk,
		Slug:          slug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IDToko:        store.ID,
		IDCategory:    req.IDCategory,
		Status:        "active",
		Berat:         req.Berat,
		SoldCount:     0,
	}

	err = u.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	// Update search index async
	go func() {
		// Simulate updating search index
		// updateSearchIndex(product.ID, product.NamaProduk)
	}()

	return product, nil
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

	// Get existing product
	product, err := u.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// Validate category exists
	_, err = u.categoryRepo.GetByID(req.IDCategory)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Check if name changed to update slug
	nameChanged := product.NamaProduk != req.NamaProduk

	// Update fields
	product.NamaProduk = req.NamaProduk
	product.HargaReseller = req.HargaReseller
	product.HargaKonsumen = req.HargaKonsumen
	product.Stok = req.Stok
	product.Deskripsi = req.Deskripsi
	product.IDCategory = req.IDCategory
	product.Berat = req.Berat

	// Update slug if name changed
	if nameChanged {
		newSlug := generateSlug(req.NamaProduk)
		
		// Check if new slug conflicts with existing products (excluding current product)
		existingProduct, err := u.productRepo.GetBySlug(newSlug)
		if err == nil && existingProduct.ID != productID {
			// Slug exists, append number
			originalSlug := newSlug
			counter := 1
			for {
				newSlug = fmt.Sprintf("%s-%d", originalSlug, counter)
				existingProduct, err := u.productRepo.GetBySlug(newSlug)
				if err != nil || existingProduct.ID == productID {
					break
				}
				counter++
			}
		}
		product.Slug = newSlug
	}

	err = u.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
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

	return u.productRepo.Delete(productID)
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

	position := len(existingPhotos) + 1

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

// Helper function to generate slug from product name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	
	// Replace spaces and special characters with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	
	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	
	return slug
}