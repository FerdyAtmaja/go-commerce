package usecase

import (
	"errors"
	"math"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"

	"gorm.io/gorm"
)

type StoreUsecase struct {
	storeRepo domain.StoreRepository
}

func NewStoreUsecase(storeRepo domain.StoreRepository) *StoreUsecase {
	return &StoreUsecase{
		storeRepo: storeRepo,
	}
}

func (u *StoreUsecase) CreateStore(userID uint, name string) (*domain.Store, error) {
	store := &domain.Store{
		UserID:      userID,
		Name:        name,
		Description: "Welcome to " + name,
	}

	if err := u.storeRepo.Create(store); err != nil {
		return nil, errors.New("failed to create store")
	}

	return store, nil
}

func (u *StoreUsecase) GetMyStore(userID uint) (*domain.Store, error) {
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, errors.New("failed to get store")
	}

	return store, nil
}

func (u *StoreUsecase) UpdateMyStore(userID uint, req *domain.UpdateStoreRequest) (*domain.Store, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, errors.New("failed to get store")
	}

	// Update store fields
	store.Name = req.Name
	store.Description = req.Description

	if err := u.storeRepo.Update(store); err != nil {
		return nil, errors.New("failed to update store")
	}

	return store, nil
}

func (u *StoreUsecase) UpdateStorePhoto(userID uint, photoURL string) (*domain.Store, error) {
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, errors.New("failed to get store")
	}

	store.PhotoURL = photoURL
	if err := u.storeRepo.Update(store); err != nil {
		return nil, errors.New("failed to update store photo")
	}

	return store, nil
}

func (u *StoreUsecase) GetStoreByID(id uint) (*domain.Store, error) {
	store, err := u.storeRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, errors.New("failed to get store")
	}

	return store, nil
}

func (u *StoreUsecase) GetAllStores(page, limit int, search string) ([]*domain.Store, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	stores, total, err := u.storeRepo.GetAll(limit, offset, search)
	if err != nil {
		return nil, response.PaginationMeta{}, errors.New("failed to get stores")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	return stores, meta, nil
}