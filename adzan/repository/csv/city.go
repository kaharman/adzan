package csv

import (
	"adzan/domain"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type csvCityRepository struct {
	baseDir string
	mu      sync.RWMutex
	loaded  bool
	cities  []domain.City
	byCode  map[int]domain.City
	byName  map[string]domain.City
}

func NewCSVCityRepository(baseDir string) domain.CityRepository {
	repo := &csvCityRepository{
		baseDir: baseDir,
		cities:  make([]domain.City, 0),
		byCode:  make(map[int]domain.City),
		byName:  make(map[string]domain.City),
	}
	_ = repo.loadCities()
	return repo
}

func (r *csvCityRepository) loadCities() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.baseDir, "cities.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	_ = header

	r.cities = make([]domain.City, 0)
	r.byCode = make(map[int]domain.City)
	r.byName = make(map[string]domain.City)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 6 {
			continue
		}

		code, _ := strconv.Atoi(record[0])
		prov := record[1]
		cityName := record[2]
		tz, _ := strconv.Atoi(record[3])
		lat, _ := strconv.ParseFloat(record[4], 32)
		lng, _ := strconv.ParseFloat(record[5], 32)

		city := domain.City{
			Code:      code,
			Province:  prov,
			City:      cityName,
			Timezone:  tz,
			Latitude:  float32(lat),
			Longitude: float32(lng),
			IDmyQuran: code,
		}

		r.cities = append(r.cities, city)
		r.byCode[code] = city
		r.byName[strings.ToUpper(cityName)] = city
	}

	r.loaded = true
	slog.Debug("Loaded cities from CSV", "total_cities", len(r.cities), "file", filePath)
	return nil
}

func (r *csvCityRepository) Search(ctx context.Context, keyword string) (res []domain.City, err error) {
	slog.Debug("CSV CityRepository Search", "keyword", keyword)
	r.mu.RLock()
	defer r.mu.RUnlock()

	kw := strings.ToUpper(strings.TrimSpace(keyword))
	for _, city := range r.cities {
		if strings.Contains(strings.ToUpper(city.City), kw) || strings.Contains(strings.ToUpper(city.Province), kw) {
			res = append(res, city)
		}
	}
	slog.Debug("CSV CityRepository Search results", "keyword", keyword, "matches", len(res))
	return res, nil
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Radius bumi dalam km
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (r *csvCityRepository) GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res domain.City, err error) {
	slog.Debug("CSV CityRepository GetByLatLong", "latitude", latitude, "longitude", longitude)
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.cities) == 0 {
		return res, errors.New("redis: nil")
	}

	minDist := math.MaxFloat64
	var closestCity domain.City
	found := false

	targetLat := float64(latitude)
	targetLon := float64(longitude)

	for _, city := range r.cities {
		dist := haversineDistance(targetLat, targetLon, float64(city.Latitude), float64(city.Longitude))
		if dist < minDist {
			minDist = dist
			closestCity = city
			found = true
		}
	}

	if !found {
		return res, errors.New("redis: nil")
	}

	return closestCity, nil
}

func (r *csvCityRepository) GetByCityCode(ctx context.Context, cityCode int) (res domain.City, err error) {
	slog.Debug("CSV CityRepository GetByCityCode", "city_code", cityCode)
	r.mu.RLock()
	defer r.mu.RUnlock()

	city, ok := r.byCode[cityCode]
	if !ok {
		return res, errors.New("redis: nil")
	}
	return city, nil
}

func (r *csvCityRepository) GetByCityName(ctx context.Context, cityName string) (res domain.City, err error) {
	slog.Debug("CSV CityRepository GetByCityName", "city_name", cityName)
	r.mu.RLock()
	defer r.mu.RUnlock()

	city, ok := r.byName[strings.ToUpper(strings.TrimSpace(cityName))]
	if !ok {
		return res, errors.New("redis: nil")
	}
	return city, nil
}

func (r *csvCityRepository) Create(ctx context.Context, city *domain.City) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cities = append(r.cities, *city)
	r.byCode[city.Code] = *city
	r.byName[strings.ToUpper(city.City)] = *city
	return nil
}

func (r *csvCityRepository) Update(ctx context.Context, city *domain.City) error {
	return r.Create(ctx, city)
}

func (r *csvCityRepository) Delete(ctx context.Context, CityCode string) error {
	code, err := strconv.Atoi(CityCode)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.byCode, code)
	for i, c := range r.cities {
		if c.Code == code {
			delete(r.byName, strings.ToUpper(c.City))
			r.cities = append(r.cities[:i], r.cities[i+1:]...)
			break
		}
	}
	return nil
}
