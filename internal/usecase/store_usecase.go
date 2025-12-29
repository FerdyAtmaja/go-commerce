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

func (u *StoreUsecase) CreateStore(userID uint64, req *domain.CreateStoreRequest) (*domain.Store, error) {
	// Check if user already has a store (including soft deleted)
	existingStore, err := u.storeRepo.GetByUserID(userID)
	if err == nil && existingStore != nil {
		return nil, errors.New("STORE_ALREADY_EXISTS")
	}

	store := &domain.Store{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		PhotoURL:    req.PhotoURL,
		Status:      "pending", // Start with pending status
		Rating:      0.0,
	}

	if err := u.storeRepo.Create(store); err != nil {
		return nil, errors.New("failed to create store")
	}

	// Get created store with relations
	return u.storeRepo.GetByID(store.ID)
}

func (u *StoreUsecase) GetMyStore(userID uint64) (*domain.Store, error) {
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("STORE_NOT_FOUND")
		}
		return nil, errors.New("failed to get store")
	}

	return store, nil
}

func (u *StoreUsecase) UpdateMyStore(userID uint64, req *domain.UpdateStoreRequest) (*domain.Store, error) {
	// Get user's store
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, errors.New("failed to get store")
	}

	// Check if store is suspended - seller cannot update suspended store
	if store.Status == "suspended" {
		return nil, errors.New("STORE_SUSPENDED_BY_ADMIN")
	}

	// Update only provided fields
	if req.Name != nil {
		store.Name = *req.Name
	}
	if req.Description != nil {
		store.Description = *req.Description
	}
	if req.PhotoURL != nil {
		store.PhotoURL = *req.PhotoURL
	}

	if err := u.storeRepo.Update(store); err != nil {
		return nil, errors.New("failed to update store")
	}

	return store, nil
}

func (u *StoreUsecase) UpdateStorePhoto(userID uint64, photoURL string) (*domain.Store, error) {
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

func (u *StoreUsecase) GetStoreByID(id uint64) (*domain.Store, error) {
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

// ActivateStore allows seller to activate their store
func (u *StoreUsecase) ActivateStore(userID uint64) error {
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status == "suspended" {
		return errors.New("STORE_SUSPENDED_BY_ADMIN")
	}

	if store.Status == "active" {
		return errors.New("STORE_ALREADY_ACTIVE")
	}

	if store.Status == "pending" {
		return errors.New("STORE_PENDING_APPROVAL")
	}

	// Check if profile is complete
	if store.Name == "" {
		return errors.New("STORE_PROFILE_INCOMPLETE")
	}

	store.Status = "active"
	return u.storeRepo.Update(store)
}

// DeactivateStore allows seller to deactivate their store
func (u *StoreUsecase) DeactivateStore(userID uint64) error {
	store, err := u.storeRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status != "active" {
		return errors.New("STORE_NOT_ACTIVE")
	}

	store.Status = "inactive"
	return u.storeRepo.Update(store)
}

// SuspendStore allows admin to suspend any store
func (u *StoreUsecase) SuspendStore(storeID uint64) error {
	store, err := u.storeRepo.GetByID(storeID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status == "suspended" {
		return errors.New("STORE_ALREADY_SUSPENDED")
	}

	store.Status = "suspended"
	return u.storeRepo.Update(store)
}

// UnsuspendStore allows admin to unsuspend a store
func (u *StoreUsecase) UnsuspendStore(storeID uint64) error {
	store, err := u.storeRepo.GetByID(storeID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status != "suspended" {
		return errors.New("STORE_NOT_SUSPENDED")
	}

	// Set to inactive, let seller activate it
	store.Status = "inactive"
	return u.storeRepo.Update(store)
}

// GetStorePublic returns store only if it's active
func (u *StoreUsecase) GetStorePublic(storeID uint64) (*domain.Store, error) {
	store, err := u.storeRepo.GetByID(storeID)
	if err != nil {
		return nil, errors.New("store not found")
	}

	if store.Status != "active" {
		return nil, errors.New("store not found")
	}

	return store, nil
}

// GetActiveStores returns only active stores for public listing
func (u *StoreUsecase) GetActiveStores(page, limit int, search string) ([]*domain.Store, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	stores, total, err := u.storeRepo.GetActiveStores(limit, offset, search)
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
// ApproveStore allows admin to approve pending store
func (u *StoreUsecase) ApproveStore(storeID uint64) error {
	store, err := u.storeRepo.GetByID(storeID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status != "pending" {
		return errors.New("STORE_NOT_PENDING")
	}

	store.Status = "active"
	return u.storeRepo.Update(store)
}

// RejectStore allows admin to reject pending store
func (u *StoreUsecase) RejectStore(storeID uint64) error {
	store, err := u.storeRepo.GetByID(storeID)
	if err != nil {
		return errors.New("store not found")
	}

	if store.Status != "pending" {
		return errors.New("STORE_NOT_PENDING")
	}

	store.Status = "inactive"
	return u.storeRepo.Update(store)
}

// GetPendingStores returns stores waiting for admin approval
func (u *StoreUsecase) GetPendingStores(page, limit int, search string) ([]*domain.Store, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	stores, total, err := u.storeRepo.GetPendingStores(limit, offset, search)
	if err != nil {
		return nil, response.PaginationMeta{}, errors.New("failed to get pending stores")
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