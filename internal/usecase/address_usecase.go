package usecase

import (
	"errors"
	"math"

	"go-commerce/internal/domain"
	"go-commerce/internal/handler/response"

	"gorm.io/gorm"
)

type AddressUsecase struct {
	addressRepo   domain.AddressRepository
	regionService domain.RegionService
}

func NewAddressUsecase(addressRepo domain.AddressRepository, regionService domain.RegionService) *AddressUsecase {
	return &AddressUsecase{
		addressRepo:   addressRepo,
		regionService: regionService,
	}
}

func (u *AddressUsecase) CreateAddress(userID uint64, req *domain.CreateAddressRequest) (*domain.Address, error) {
	// Validate province and city IDs with Indonesia region data
	if err := u.regionService.ValidateProvinceAndCity(req.ProvinceID, req.CityID); err != nil {
		return nil, errors.New("invalid province or city: " + err.Error())
	}
	
	address := &domain.Address{
		UserID:     userID,
		Name:       req.Name,
		Detail:     req.Detail,
		Phone:      req.Phone,
		ProvinceID: req.ProvinceID,
		CityID:     req.CityID,
		PostalCode: req.PostalCode,
	}

	if err := u.addressRepo.Create(address); err != nil {
		return nil, errors.New("failed to create address")
	}

	return address, nil
}

func (u *AddressUsecase) GetAddressByID(addressID, userID uint64) (*domain.Address, error) {
	// Check ownership
	if !u.addressRepo.CheckOwnership(addressID, userID) {
		return nil, errors.New("address not found or access denied")
	}

	address, err := u.addressRepo.GetByID(addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, errors.New("failed to get address")
	}

	return address, nil
}

func (u *AddressUsecase) GetMyAddresses(userID uint64, page, limit int) ([]*domain.Address, response.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	addresses, total, err := u.addressRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, response.PaginationMeta{}, errors.New("failed to get addresses")
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	meta := response.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
	}

	return addresses, meta, nil
}

func (u *AddressUsecase) UpdateAddress(addressID, userID uint64, req *domain.UpdateAddressRequest) (*domain.Address, error) {
	// Check ownership
	if !u.addressRepo.CheckOwnership(addressID, userID) {
		return nil, errors.New("address not found or access denied")
	}

	address, err := u.addressRepo.GetByID(addressID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, errors.New("failed to get address")
	}

	// Validate province and city IDs with Indonesia region data
	if err := u.regionService.ValidateProvinceAndCity(req.ProvinceID, req.CityID); err != nil {
		return nil, errors.New("invalid province or city: " + err.Error())
	}

	// Update address fields
	address.Name = req.Name
	address.Detail = req.Detail
	address.Phone = req.Phone
	address.ProvinceID = req.ProvinceID
	address.CityID = req.CityID
	address.PostalCode = req.PostalCode

	if err := u.addressRepo.Update(address); err != nil {
		return nil, errors.New("failed to update address")
	}

	return address, nil
}

func (u *AddressUsecase) DeleteAddress(addressID, userID uint64) error {
	// Check ownership
	if !u.addressRepo.CheckOwnership(addressID, userID) {
		return errors.New("address not found or access denied")
	}

	if err := u.addressRepo.Delete(addressID); err != nil {
		return errors.New("failed to delete address")
	}

	return nil
}

func (u *AddressUsecase) GetProvinces() ([]*domain.Province, error) {
	provinces, err := u.regionService.GetProvinces()
	if err != nil {
		return nil, errors.New("failed to get provinces: " + err.Error())
	}
	return provinces, nil
}

func (u *AddressUsecase) GetCitiesByProvince(provinceID string) ([]*domain.City, error) {
	cities, err := u.regionService.GetCitiesByProvince(provinceID)
	if err != nil {
		return nil, errors.New("failed to get cities: " + err.Error())
	}
	return cities, nil
}