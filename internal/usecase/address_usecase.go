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
	address := &domain.Address{
		UserID:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
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

	// Update address fields
	address.JudulAlamat = req.JudulAlamat
	address.NamaPenerima = req.NamaPenerima
	address.NoTelp = req.NoTelp
	address.DetailAlamat = req.DetailAlamat
	address.KodePos = req.KodePos
	address.IsDefault = req.IsDefault

	// Handle default address logic
	if req.IsDefault && !address.IsDefault {
		// Setting as new default
		if err := u.addressRepo.SetDefault(addressID, userID); err != nil {
			return nil, errors.New("failed to set default address")
		}
	}

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