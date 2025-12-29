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
	// Validate required fields are not empty
	if len(req.JudulAlamat) < 2 || len(req.JudulAlamat) > 255 {
		return nil, errors.New("judul_alamat must be between 2 and 255 characters")
	}
	if len(req.NamaPenerima) < 2 || len(req.NamaPenerima) > 255 {
		return nil, errors.New("nama_penerima must be between 2 and 255 characters")
	}
	if len(req.DetailAlamat) < 2 {
		return nil, errors.New("detail_alamat must be at least 2 characters")
	}
	if req.ProvinceID == "" {
		return nil, errors.New("province_id is required")
	}
	if req.CityID == "" {
		return nil, errors.New("city_id is required")
	}
	if len(req.NoTelp) > 0 && (len(req.NoTelp) < 10 || len(req.NoTelp) > 20) {
		return nil, errors.New("notelp must be between 10 and 20 characters when provided")
	}
	if len(req.KodePos) > 10 {
		return nil, errors.New("kode_pos must be maximum 10 characters")
	}

	// Validate province and city using region service
	if err := u.regionService.ValidateProvinceAndCity(req.ProvinceID, req.CityID); err != nil {
		return nil, errors.New("invalid province or city: " + err.Error())
	}

	address := &domain.Address{
		UserID:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
		ProvinceID:   req.ProvinceID,
		CityID:       req.CityID,
		KodePos:      req.KodePos,
		IsDefault:    req.IsDefault,
	}

	// If this is set as default, unset other default addresses
	if req.IsDefault {
		if err := u.addressRepo.SetDefault(0, userID); err != nil {
			// Continue even if this fails
		}
	}

	if err := u.addressRepo.Create(address); err != nil {
		return nil, errors.New("failed to create address")
	}

	// Set as default after creation if requested
	if req.IsDefault {
		if err := u.addressRepo.SetDefault(address.ID, userID); err != nil {
			// Log error but don't fail the creation
		}
	}

	// Populate province and city names for response
	u.populateRegionNames(address)

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

	// Populate province and city names
	u.populateRegionNames(address)

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

	// Populate province and city names for all addresses
	u.populateRegionNamesForAddresses(addresses)

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

	// Validate province and city if provided
	if req.ProvinceID != nil && req.CityID != nil {
		if err := u.regionService.ValidateProvinceAndCity(*req.ProvinceID, *req.CityID); err != nil {
			return nil, errors.New("invalid province or city: " + err.Error())
		}
	}

	// Update only provided fields with validation
	if req.JudulAlamat != nil {
		if len(*req.JudulAlamat) < 2 || len(*req.JudulAlamat) > 255 {
			return nil, errors.New("judul_alamat must be between 2 and 255 characters")
		}
		address.JudulAlamat = *req.JudulAlamat
	}
	if req.NamaPenerima != nil {
		if len(*req.NamaPenerima) < 2 || len(*req.NamaPenerima) > 255 {
			return nil, errors.New("nama_penerima must be between 2 and 255 characters")
		}
		address.NamaPenerima = *req.NamaPenerima
	}
	if req.NoTelp != nil {
		if len(*req.NoTelp) > 0 && (len(*req.NoTelp) < 10 || len(*req.NoTelp) > 20) {
			return nil, errors.New("notelp must be between 10 and 20 characters when provided")
		}
		address.NoTelp = *req.NoTelp
	}
	if req.DetailAlamat != nil {
		if len(*req.DetailAlamat) < 2 {
			return nil, errors.New("detail_alamat must be at least 2 characters")
		}
		address.DetailAlamat = *req.DetailAlamat
	}
	if req.ProvinceID != nil {
		address.ProvinceID = *req.ProvinceID
	}
	if req.CityID != nil {
		address.CityID = *req.CityID
	}
	if req.KodePos != nil {
		if len(*req.KodePos) > 10 {
			return nil, errors.New("kode_pos must be maximum 10 characters")
		}
		address.KodePos = *req.KodePos
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	// Handle default address logic
	if req.IsDefault != nil && *req.IsDefault && !address.IsDefault {
		// Setting as new default
		if err := u.addressRepo.SetDefault(addressID, userID); err != nil {
			return nil, errors.New("failed to set default address")
		}
	}

	if err := u.addressRepo.Update(address); err != nil {
		return nil, errors.New("failed to update address")
	}

	// Populate province and city names for response
	u.populateRegionNames(address)

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

func (u *AddressUsecase) SetDefaultAddress(addressID, userID uint64) error {
	// Check ownership
	if !u.addressRepo.CheckOwnership(addressID, userID) {
		return errors.New("address not found or access denied")
	}

	if err := u.addressRepo.SetDefault(addressID, userID); err != nil {
		return errors.New("failed to set default address")
	}

	return nil
}

func (u *AddressUsecase) GetDefaultAddress(userID uint64) (*domain.Address, error) {
	address, err := u.addressRepo.GetDefaultByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no default address found")
		}
		return nil, errors.New("failed to get default address")
	}

	// Populate province and city names
	u.populateRegionNames(address)

	return address, nil
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

// populateRegionNames fills province and city names for display
func (u *AddressUsecase) populateRegionNames(address *domain.Address) {
	// Get province name
	if provinces, err := u.regionService.GetProvinces(); err == nil {
		for _, province := range provinces {
			if province.ID == address.ProvinceID {
				address.ProvinceName = province.Name
				break
			}
		}
	}

	// Get city name
	if cities, err := u.regionService.GetCitiesByProvince(address.ProvinceID); err == nil {
		for _, city := range cities {
			if city.ID == address.CityID {
				address.CityName = city.Name
				break
			}
		}
	}
}

// populateRegionNamesForAddresses fills province and city names for multiple addresses
func (u *AddressUsecase) populateRegionNamesForAddresses(addresses []*domain.Address) {
	for _, address := range addresses {
		u.populateRegionNames(address)
	}
}