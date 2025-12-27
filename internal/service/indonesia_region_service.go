package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-commerce/internal/domain"
)

type IndonesiaRegionService struct {
	baseURL string
	client  *http.Client
}

func NewIndonesiaRegionService() domain.RegionService {
	return &IndonesiaRegionService{
		baseURL: "https://emsifa.github.io/api-wilayah-indonesia/api",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *IndonesiaRegionService) GetProvinces() ([]*domain.Province, error) {
	url := fmt.Sprintf("%s/provinces.json", s.baseURL)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var provinces []*domain.Province
	if err := json.Unmarshal(body, &provinces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provinces: %w", err)
	}

	return provinces, nil
}

func (s *IndonesiaRegionService) GetCitiesByProvince(provinceID string) ([]*domain.City, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", s.baseURL, provinceID)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var cities []*domain.City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cities: %w", err)
	}

	return cities, nil
}

func (s *IndonesiaRegionService) ValidateProvinceAndCity(provinceID, cityID string) error {
	// Validate province exists
	url := fmt.Sprintf("%s/province/%s.json", s.baseURL, provinceID)
	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to validate province: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid province ID: %s", provinceID)
	}

	// Validate city exists and belongs to province
	url = fmt.Sprintf("%s/regency/%s.json", s.baseURL, cityID)
	resp, err = s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to validate city: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid city ID: %s", cityID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read city response: %w", err)
	}

	var city domain.City
	if err := json.Unmarshal(body, &city); err != nil {
		return fmt.Errorf("failed to unmarshal city: %w", err)
	}

	if city.ProvinceID != provinceID {
		return fmt.Errorf("city %s does not belong to province %s", cityID, provinceID)
	}

	return nil
}