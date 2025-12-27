package usecase

import (
	"errors"
	"testing"

	"go-commerce/internal/domain"
	"go-commerce/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductUsecase_CreateProduct_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	store := &domain.Store{
		ID:     1,
		UserID: userID,
		Name:   "Test Store",
	}
	category := &domain.Category{
		ID:   1,
		Name: "Electronics",
	}
	req := &domain.CreateProductRequest{
		NamaProduk:    "iPhone 15",
		HargaReseller: 15000000,
		HargaKonsumen: 17000000,
		Stok:          10,
		Deskripsi:     "Latest iPhone",
		IDCategory:    1,
		Berat:         200,
	}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	categoryRepo.On("GetByID", req.IDCategory).Return(category, nil)
	productRepo.On("GetBySlug", "iphone-15").Return(nil, errors.New("not found"))
	productRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

	// Execute
	product, err := usecase.CreateProduct(userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.NamaProduk, product.NamaProduk)
	assert.Equal(t, "iphone-15", product.Slug)
	assert.Equal(t, store.ID, product.IDToko)
	assert.Equal(t, req.IDCategory, product.IDCategory)
	assert.Equal(t, "active", product.Status)

	// Verify expectations
	storeRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_StoreNotFound(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	req := &domain.CreateProductRequest{
		NamaProduk:    "iPhone 15",
		HargaReseller: 15000000,
		HargaKonsumen: 17000000,
		IDCategory:    1,
	}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(nil, errors.New("store not found"))

	// Execute
	product, err := usecase.CreateProduct(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "store not found", err.Error())

	storeRepo.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_CategoryNotFound(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}
	req := &domain.CreateProductRequest{
		NamaProduk:    "iPhone 15",
		HargaReseller: 15000000,
		HargaKonsumen: 17000000,
		IDCategory:    999,
	}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	categoryRepo.On("GetByID", req.IDCategory).Return(nil, errors.New("category not found"))

	// Execute
	product, err := usecase.CreateProduct(userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "category not found", err.Error())

	storeRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestProductUsecase_GetMyProducts_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}
	products := []*domain.Product{
		{ID: 1, NamaProduk: "Product 1", IDToko: store.ID},
		{ID: 2, NamaProduk: "Product 2", IDToko: store.ID},
	}
	total := int64(2)

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	productRepo.On("GetByTokoID", store.ID, 10, 0, "").Return(products, total, nil)

	// Execute
	result, resultTotal, err := usecase.GetMyProducts(userID, 1, 10, "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, products, result)
	assert.Equal(t, total, resultTotal)

	storeRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	productID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}
	category := &domain.Category{ID: 1, Name: "Electronics"}
	existingProduct := &domain.Product{
		ID:         productID,
		NamaProduk: "Old Product",
		Slug:       "old-product",
		IDToko:     store.ID,
		IDCategory: 1,
	}
	req := &domain.UpdateProductRequest{
		NamaProduk:    "Updated Product",
		HargaReseller: 20000000,
		HargaKonsumen: 22000000,
		Stok:          5,
		IDCategory:    1,
		Berat:         300,
	}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	productRepo.On("CheckOwnership", productID, store.ID).Return(nil)
	productRepo.On("GetByID", productID).Return(existingProduct, nil)
	categoryRepo.On("GetByID", req.IDCategory).Return(category, nil)
	productRepo.On("GetBySlug", "updated-product").Return(nil, errors.New("not found"))
	productRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(nil)

	// Execute
	product, err := usecase.UpdateProduct(userID, productID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.NamaProduk, product.NamaProduk)
	assert.Equal(t, "updated-product", product.Slug)

	// Verify expectations
	storeRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_AccessDenied(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	productID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}
	req := &domain.UpdateProductRequest{
		NamaProduk: "Updated Product",
		IDCategory: 1,
	}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	productRepo.On("CheckOwnership", productID, store.ID).Return(errors.New("access denied: product not owned by this store"))

	// Execute
	product, err := usecase.UpdateProduct(userID, productID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "access denied: product not owned by this store", err.Error())

	storeRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	productID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	productRepo.On("CheckOwnership", productID, store.ID).Return(nil)
	productRepo.On("Delete", productID).Return(nil)

	// Execute
	err := usecase.DeleteProduct(userID, productID)

	// Assert
	assert.NoError(t, err)

	storeRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductUsecase_AddProductPhoto_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	userID := uint64(1)
	productID := uint64(1)
	store := &domain.Store{ID: 1, UserID: userID}
	photoURL := "/uploads/products/test.jpg"
	isPrimary := true
	existingPhotos := []*domain.PhotoProduk{}

	// Setup expectations
	storeRepo.On("GetByUserID", userID).Return(store, nil)
	productRepo.On("CheckOwnership", productID, store.ID).Return(nil)
	photoRepo.On("GetByProductID", productID).Return(existingPhotos, nil)
	photoRepo.On("Create", mock.AnythingOfType("*domain.PhotoProduk")).Return(nil).Run(func(args mock.Arguments) {
		photo := args.Get(0).(*domain.PhotoProduk)
		photo.ID = 1 // Simulate database ID assignment
	})
	photoRepo.On("SetPrimary", productID, uint64(1)).Return(nil)

	// Execute
	photo, err := usecase.AddProductPhoto(userID, productID, photoURL, isPrimary)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, photo)
	assert.Equal(t, productID, photo.IDProduk)
	assert.Equal(t, photoURL, photo.URL)
	assert.Equal(t, isPrimary, photo.IsPrimary)
	assert.Equal(t, 1, photo.Position)

	// Verify expectations
	storeRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
	photoRepo.AssertExpectations(t)
}

func TestProductUsecase_GetAllProducts_Success(t *testing.T) {
	// Setup mocks
	productRepo := new(mocks.ProductRepositoryMock)
	photoRepo := new(mocks.PhotoProdukRepositoryMock)
	storeRepo := new(mocks.StoreRepositoryMock)
	categoryRepo := new(mocks.CategoryRepositoryMock)

	usecase := NewProductUsecase(productRepo, photoRepo, storeRepo, categoryRepo)

	// Test data
	filter := &domain.ProductFilter{
		Search:     "iPhone",
		CategoryID: "1",
		Page:       1,
		Limit:      10,
	}
	products := []*domain.Product{
		{ID: 1, NamaProduk: "iPhone 15", Status: "active"},
		{ID: 2, NamaProduk: "iPhone 14", Status: "active"},
	}
	total := int64(2)

	// Setup expectations
	productRepo.On("GetAll", 10, 0, "iPhone", "1").Return(products, total, nil)

	// Execute
	result, resultTotal, err := usecase.GetAllProducts(filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, products, result)
	assert.Equal(t, total, resultTotal)

	productRepo.AssertExpectations(t)
}