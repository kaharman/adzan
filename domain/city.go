package domain

import "context"

type City struct {
	Code      int     `db:"code"`
	Province  string  `db:"province"`
	City      string  `db:"city"`
	Latitude  float32 `db:"latitude"`
	Longitude float32 `db:"longitude"`
	Timezone  int     `db:"timezone"`
	IDmyQuran int     `db:"id_my_quran"`
}

type CityResponseDTO struct {
	CityCode  int     `json:"city_code"`
	City      string  `json:"city"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timezone  int     `json:"timezone"`
}

type CityUsecase interface {
	GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res CityResponseDTO, err error)
	GetByCityCode(ctx context.Context, cityCode int) (res CityResponseDTO, err error)
	Search(ctx context.Context, keyword string) (res []CityResponseDTO, err error)
}

type CityRepository interface {
	Search(ctx context.Context, keyword string) (res []City, err error)
	GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res City, err error)
	GetByCityCode(ctx context.Context, cityCode int) (res City, err error)
	GetByCityName(ctx context.Context, cityName string) (res City, err error)
	Create(ctx context.Context, city *City) (err error)
	Update(ctx context.Context, city *City) (err error)
	Delete(ctx context.Context, CityCode string) (err error)
}
