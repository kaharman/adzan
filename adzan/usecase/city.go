package usecase

import (
	"adzan/domain"
	"context"
)

type cityUsecase struct {
	cityRepo domain.CityRepository
}

func NewCityUsecase(cityRepo domain.CityRepository) domain.CityUsecase {
	return &cityUsecase{
		cityRepo: cityRepo,
	}
}

func (c *cityUsecase) GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res domain.CityResponseDTO, err error) {
	city, err := c.cityRepo.GetByLatLong(ctx, latitude, longitude)
	if err != nil {
		return
	}
	res = domain.CityResponseDTO{
		CityCode:  city.Code,
		City:      city.City,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Timezone:  city.Timezone,
	}

	return
}

func (c *cityUsecase) GetByCityCode(ctx context.Context, cityCode int) (res domain.CityResponseDTO, err error) {
	city, err := c.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}
	res = domain.CityResponseDTO{
		CityCode:  city.Code,
		City:      city.City,
		Latitude:  city.Latitude,
		Longitude: city.Longitude,
		Timezone:  city.Timezone,
	}

	return
}

func (c *cityUsecase) Search(ctx context.Context, keyword string) (res []domain.CityResponseDTO, err error) {
	cities, err := c.cityRepo.Search(ctx, keyword)
	if err != nil {
		return
	}
	for _, city := range cities {
		res = append(res, domain.CityResponseDTO{
			CityCode:  city.Code,
			City:      city.City,
			Latitude:  city.Latitude,
			Longitude: city.Longitude,
			Timezone:  city.Timezone,
		})
	}

	return
}
