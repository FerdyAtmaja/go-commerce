package domain

type RegionService interface {
	GetProvinces() ([]*Province, error)
	GetCitiesByProvince(provinceID string) ([]*City, error)
	ValidateProvinceAndCity(provinceID, cityID string) error
}