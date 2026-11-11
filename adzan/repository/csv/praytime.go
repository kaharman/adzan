package csv

import (
	"adzan/domain"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type csvPrayTimeRepository struct {
	baseDir string
	mu      sync.RWMutex
	cache   map[int]map[int]map[int]domain.PrayTimeSchedule
}

func NewCSVPrayTimeRepository(baseDir string) domain.PrayTimeRepository {
	return &csvPrayTimeRepository{
		baseDir: baseDir,
		cache:   make(map[int]map[int]map[int]domain.PrayTimeSchedule),
	}
}

func (r *csvPrayTimeRepository) loadCityFile(cityCode int) (map[int]map[int]domain.PrayTimeSchedule, error) {
	r.mu.RLock()
	if data, ok := r.cache[cityCode]; ok {
		r.mu.RUnlock()
		slog.Debug("CSV pray time cache hit", "city_code", cityCode)
		return data, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	if data, ok := r.cache[cityCode]; ok {
		return data, nil
	}

	filePath := filepath.Join(r.baseDir, "praytimes", fmt.Sprintf("%d.csv", cityCode))
	slog.Debug("Loading CSV pray times file from disk", "city_code", cityCode, "path", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("redis: nil") // Selaras dengan error handling di HTTP controller
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	_ = header

	cityData := make(map[int]map[int]domain.PrayTimeSchedule)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 10 {
			continue
		}

		month, errM := strconv.Atoi(record[0])
		day, errD := strconv.Atoi(record[1])
		if errM != nil || errD != nil {
			continue
		}

		if _, ok := cityData[month]; !ok {
			cityData[month] = make(map[int]domain.PrayTimeSchedule)
		}

		cityData[month][day] = domain.PrayTimeSchedule{
			Imsak:    record[2],
			Fajr:     record[3],
			Sunrise:  record[4],
			Dhuhr:    record[5],
			Asr:      record[6],
			Maghrib:  record[7],
			Isha:     record[8],
			Midnight: record[9],
		}
	}

	r.cache[cityCode] = cityData
	return cityData, nil
}

func (r *csvPrayTimeRepository) GetPrayTime(ctx context.Context, cityCode int, month int, day int) (res domain.PrayTimeSchedule, err error) {
	slog.Debug("CSV GetPrayTime query", "city_code", cityCode, "month", month, "day", day)
	cityData, err := r.loadCityFile(cityCode)
	if err != nil {
		return
	}

	monthData, ok := cityData[month]
	if !ok {
		return res, errors.New("redis: nil")
	}

	schedule, ok := monthData[day]
	if !ok {
		return res, errors.New("redis: nil")
	}

	return schedule, nil
}

func (r *csvPrayTimeRepository) GetMonthlyPrayTime(ctx context.Context, cityCode int, month int) (res map[int]domain.PrayTimeSchedule, err error) {
	cityData, err := r.loadCityFile(cityCode)
	if err != nil {
		return
	}

	monthData, ok := cityData[month]
	if !ok || len(monthData) == 0 {
		return nil, errors.New("redis: nil")
	}

	return monthData, nil
}

func (r *csvPrayTimeRepository) SaveDaily(ctx context.Context, cityCode int, month int, day int, schedule *domain.PrayTimeSchedule) error {
	cityData, _ := r.loadCityFile(cityCode)
	r.mu.Lock()
	defer r.mu.Unlock()

	if cityData == nil {
		cityData = make(map[int]map[int]domain.PrayTimeSchedule)
		r.cache[cityCode] = cityData
	}
	if _, ok := cityData[month]; !ok {
		cityData[month] = make(map[int]domain.PrayTimeSchedule)
	}
	cityData[month][day] = *schedule
	return nil
}

func (r *csvPrayTimeRepository) SaveMonthly(ctx context.Context, cityCode int, month int, schedules map[int]domain.PrayTimeSchedule) error {
	cityData, _ := r.loadCityFile(cityCode)
	r.mu.Lock()
	defer r.mu.Unlock()

	if cityData == nil {
		cityData = make(map[int]map[int]domain.PrayTimeSchedule)
		r.cache[cityCode] = cityData
	}
	cityData[month] = schedules
	return nil
}

func (r *csvPrayTimeRepository) Delete(ctx context.Context, cityCode int, month int, day int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cityData, ok := r.cache[cityCode]; ok {
		if monthData, ok := cityData[month]; ok {
			delete(monthData, day)
		}
	}
	return nil
}
