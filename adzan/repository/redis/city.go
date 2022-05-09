package redis

import (
	"adzan/domain"
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type city struct {
	Province  string `redis:"province"`
	City      string `redis:"city"`
	Timezone  int    `redis:"timezone"`
	IDmyQuran int    `redis:"id_my_quran"`
}

type redisCityRepository struct {
	Conn *redis.Client
}

func NewRedisCityRepository(Conn *redis.Client) domain.CityRepository {
	return &redisCityRepository{Conn}
}

func (r *redisCityRepository) Search(ctx context.Context, offset, limit int64, sort, province string, desc string) (res []domain.City, err error) {
	return
}

func (r *redisCityRepository) GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res domain.City, err error) {
	location, err := r.Conn.GeoSearchLocation(ctx, "city_locations", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Latitude:   float64(latitude),
			Longitude:  float64(longitude),
			Radius:     100,
			RadiusUnit: "km",
			Sort:       "ASC",
			Count:      5,
		},
		WithCoord: true,
	}).Result()
	if err != nil {
		return
	}
	var city city
	err = r.Conn.HGetAll(ctx, "cities:"+location[0].Name).Scan(&city)
	if err != nil {
		return
	}
	if city.Timezone == 0 {
		err = errors.New("redis: nil")
		return
	}
	cityCode, err := strconv.ParseInt(location[0].Name, 10, 64)
	if err != nil {
		return
	}
	res = domain.City{
		Code:      int(cityCode),
		Province:  city.Province,
		City:      city.City,
		Latitude:  float32(location[0].Latitude),
		Longitude: float32(location[0].Longitude),
		Timezone:  city.Timezone,
		IDmyQuran: city.IDmyQuran,
	}

	return
}

func (r *redisCityRepository) GetByCityCode(ctx context.Context, cityCode int) (res domain.City, err error) {
	var city city
	err = r.Conn.HGetAll(ctx, "cities:"+strconv.FormatInt(int64(cityCode), 10)).Scan(&city)
	if err != nil {
		return
	}
	if city.Timezone == 0 {
		err = errors.New("redis: nil")
		return
	}
	pos, err := r.Conn.GeoPos(ctx, "city_locations", strconv.FormatInt(int64(cityCode), 10)).Result()
	if err != nil {
		return
	}
	res = domain.City{
		Code:      cityCode,
		Province:  city.Province,
		City:      city.City,
		Latitude:  float32(pos[0].Latitude),
		Longitude: float32(pos[0].Longitude),
		Timezone:  city.Timezone,
		IDmyQuran: city.IDmyQuran,
	}

	return
}

func (r *redisCityRepository) GetByCityName(ctx context.Context, cityName string) (res domain.City, err error) {
	var city city
	cityCode, err := r.Conn.ZScore(ctx, "cities", cityName).Result()
	if err != nil {
		return
	}
	err = r.Conn.HGetAll(ctx, "cities:"+strconv.FormatInt(int64(cityCode), 10)).Scan(&city)
	if err != nil {
		return
	}
	if city.Timezone == 0 {
		err = errors.New("redis: nil")
		return
	}
	pos, err := r.Conn.GeoPos(ctx, "city_locations", strconv.FormatInt(int64(cityCode), 10)).Result()
	if err != nil {
		return
	}
	res = domain.City{
		Code:      int(cityCode),
		Province:  city.Province,
		City:      city.City,
		Latitude:  float32(pos[0].Latitude),
		Longitude: float32(pos[0].Longitude),
		Timezone:  city.Timezone,
		IDmyQuran: city.IDmyQuran,
	}

	return
}

func (r *redisCityRepository) Create(ctx context.Context, city *domain.City) (err error) {
	return
}

func (r *redisCityRepository) Update(ctx context.Context, city *domain.City) (err error) {
	return
}

func (r *redisCityRepository) Delete(ctx context.Context, CityCode string) (err error) {
	return
}
