package mysql

import (
	"go-commerce/internal/domain"
	"log"

	"gorm.io/gorm"
)

type productLogRepository struct {
	db *gorm.DB
}

func NewProductLogRepository(db *gorm.DB) domain.ProductLogRepository {
	return &productLogRepository{db: db}
}

func (r *productLogRepository) Create(productLog *domain.ProductLog) error {
	return r.db.Create(productLog).Error
}

func (r *productLogRepository) CreateAsync(productLog *domain.ProductLog) {
	go func() {
		if err := r.db.Create(productLog).Error; err != nil {
			log.Printf("Failed to create product log asynchronously: %v", err)
		}
	}()
}

func (r *productLogRepository) GetByID(id uint64) (*domain.ProductLog, error) {
	var productLog domain.ProductLog
	err := r.db.First(&productLog, id).Error
	if err != nil {
		return nil, err
	}
	return &productLog, nil
}

func (r *productLogRepository) GetByProductID(productID uint64) ([]*domain.ProductLog, error) {
	var logs []*domain.ProductLog
	err := r.db.Where("id_produk = ?", productID).Find(&logs).Error
	return logs, err
}