package mysql

import (
	"errors"
	"strings"

	"go-commerce/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("Photos").
		Where("id = ?", id).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Toko").Preload("Category").Preload("Photos").
		Where("slug = ?", slug).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetByTokoID(tokoID uint64, limit, offset int, search string) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).Where("id_toko = ?", tokoID)

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(nama_produk) LIKE ? OR LOWER(deskripsi) LIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	err := query.Preload("Category").Preload("Photos").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetAll(limit, offset int, search, categoryID string) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).Where("status = ?", "active")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(nama_produk) LIKE ? OR LOWER(deskripsi) LIKE ?", searchPattern, searchPattern)
	}

	if categoryID != "" {
		query = query.Where("id_category = ?", categoryID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	err := query.Preload("Toko").Preload("Category").Preload("Photos").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) CheckOwnership(productID, tokoID uint64) error {
	var count int64
	err := r.db.Model(&domain.Product{}).
		Where("id = ? AND id_toko = ?", productID, tokoID).
		Count(&count).Error
	
	if err != nil {
		return err
	}
	
	if count == 0 {
		return errors.New("access denied: product not owned by this store")
	}
	
	return nil
}

// PhotoProduk Repository
type photoProdukRepository struct {
	db *gorm.DB
}

func NewPhotoProdukRepository(db *gorm.DB) domain.PhotoProdukRepository {
	return &photoProdukRepository{db: db}
}

func (r *photoProdukRepository) Create(photo *domain.PhotoProduk) error {
	return r.db.Create(photo).Error
}

func (r *photoProdukRepository) GetByProductID(productID uint64) ([]*domain.PhotoProduk, error) {
	var photos []*domain.PhotoProduk
	err := r.db.Where("id_produk = ?", productID).
		Order("is_primary DESC, position ASC").
		Find(&photos).Error
	return photos, err
}

func (r *photoProdukRepository) Update(photo *domain.PhotoProduk) error {
	return r.db.Save(photo).Error
}

func (r *photoProdukRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.PhotoProduk{}, id).Error
}

func (r *photoProdukRepository) SetPrimary(productID, photoID uint64) error {
	tx := r.db.Begin()
	
	// Reset all photos to non-primary
	if err := tx.Model(&domain.PhotoProduk{}).
		Where("id_produk = ?", productID).
		Update("is_primary", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Set selected photo as primary
	if err := tx.Model(&domain.PhotoProduk{}).
		Where("id = ? AND id_produk = ?", photoID, productID).
		Update("is_primary", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}