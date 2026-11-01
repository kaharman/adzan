package usecase

import (
	"adzan/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type apiMyQuranKoordinat struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Lintang string  `json:"lintang"`
	Bujur   string  `json:"bujur"`
}

type apiMyQuranJadwal struct {
	Tanggal string `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Terbit  string `json:"terbit"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
	Date    string `json:"date"`
}

type apiMyQuranData struct {
	ID        string              `json:"id"`
	Lokasi    string              `json:"lokasi"`
	Daerah    string              `json:"daerah"`
	Koordinat apiMyQuranKoordinat `json:"koordinat"`
	Jadwal    []apiMyQuranJadwal  `json:"jadwal"`
}

type apiMyQuranMonthlyResponse struct {
	Status bool           `json:"status"`
	Data   apiMyQuranData `json:"data"`
}

type prayTimeUsecase struct {
	cityRepo     domain.CityRepository
	prayTimeRepo domain.PrayTimeRepository
}

func NewPrayTimeUsecase(cityRepo domain.CityRepository, prayTimeRepo domain.PrayTimeRepository) domain.PrayTimeUsecase {
	return &prayTimeUsecase{
		cityRepo:     cityRepo,
		prayTimeRepo: prayTimeRepo,
	}
}

func getTimeLocation(timezone int) *time.Location {
	switch timezone {
	case 7:
		return time.FixedZone("UTC+7", 7*60*60)
	case 8:
		return time.FixedZone("UTC+8", 8*60*60)
	case 9:
		return time.FixedZone("UTC+9", 9*60*60)
	default:
		return time.FixedZone("UTC+7", 7*60*60)
	}
}

func parsePrayerTimeToUnix(year int, month time.Month, day int, timeStr string, loc *time.Location, isMidnight bool, maghribHour int) int {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0
	}
	targetDay := day
	if isMidnight && t.Hour() < maghribHour {
		targetDay = day + 1
	}
	dt := time.Date(year, month, targetDay, t.Hour(), t.Minute(), 0, 0, loc)
	return int(dt.Unix())
}

func buildTimeResponseDTO(year int, month time.Month, day int, schedule domain.PrayTimeSchedule, loc *time.Location) domain.TimeResponseDTO {
	maghribTime, _ := time.Parse("15:04", schedule.Maghrib)
	maghribHour := maghribTime.Hour()

	return domain.TimeResponseDTO{
		Imsak:    parsePrayerTimeToUnix(year, month, day, schedule.Imsak, loc, false, 0),
		Fajr:     parsePrayerTimeToUnix(year, month, day, schedule.Fajr, loc, false, 0),
		Sunrise:  parsePrayerTimeToUnix(year, month, day, schedule.Sunrise, loc, false, 0),
		Dhuhr:    parsePrayerTimeToUnix(year, month, day, schedule.Dhuhr, loc, false, 0),
		Asr:      parsePrayerTimeToUnix(year, month, day, schedule.Asr, loc, false, 0),
		Maghrib:  parsePrayerTimeToUnix(year, month, day, schedule.Maghrib, loc, false, 0),
		Isha:     parsePrayerTimeToUnix(year, month, day, schedule.Isha, loc, false, 0),
		Midnight: parsePrayerTimeToUnix(year, month, day, schedule.Midnight, loc, true, maghribHour),
	}
}

func (p *prayTimeUsecase) getDaily(ctx context.Context, city domain.City, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	datetime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return
	}

	schedule, err := p.prayTimeRepo.GetPrayTime(ctx, city.Code, int(datetime.Month()), datetime.Day())
	if err != nil {
		return
	}
	var day string
	switch lang {
	case "en":
		day = domain.DayEn[datetime.Weekday()]
	case "id":
		day = domain.DayId[datetime.Weekday()]
	case "ar":
		day = domain.DayAr[datetime.Weekday()]
	}

	loc := getTimeLocation(city.Timezone)
	res = domain.PrayTimeResponseDTO{
		Location: domain.CityResponseDTO{
			CityCode:  city.Code,
			City:      city.City,
			Latitude:  city.Latitude,
			Longitude: city.Longitude,
			Timezone:  city.Timezone,
		},
		Date: domain.DateResponseDTO{
			Day:  day,
			Date: date,
		},
		Times: buildTimeResponseDTO(datetime.Year(), datetime.Month(), datetime.Day(), schedule, loc),
	}

	return
}

func (p *prayTimeUsecase) getMonthly(ctx context.Context, city domain.City, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	datetime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	dayOfMonth := datetime.AddDate(0, 1, -1).Day()

	schedules, errBlock := p.prayTimeRepo.GetMonthlyPrayTime(ctx, city.Code, month)
	if errBlock != nil {
		err = errBlock
		return
	}

	loc := getTimeLocation(city.Timezone)

	for i := 0; i < dayOfMonth; i++ {
		dayNum := datetime.Day()
		schedule, ok := schedules[dayNum]
		if !ok {
			datetime = datetime.AddDate(0, 0, 1)
			continue
		}

		var day string
		switch lang {
		case "en":
			day = domain.DayEn[datetime.Weekday()]
		case "id":
			day = domain.DayId[datetime.Weekday()]
		case "ar":
			day = domain.DayAr[datetime.Weekday()]
		}
		result := domain.PrayTimeResponseDTO{
			Location: domain.CityResponseDTO{
				CityCode:  city.Code,
				City:      city.City,
				Latitude:  city.Latitude,
				Longitude: city.Longitude,
				Timezone:  city.Timezone,
			},
			Date: domain.DateResponseDTO{
				Day:  day,
				Date: datetime.Format("2006-01-02"),
			},
			Times: buildTimeResponseDTO(year, time.Month(month), dayNum, schedule, loc),
		}
		res = append(res, result)
		datetime = datetime.AddDate(0, 0, 1)
	}

	return
}

func (p *prayTimeUsecase) getYearly(ctx context.Context, city domain.City, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	loc := getTimeLocation(city.Timezone)

	for m := 1; m <= 12; m++ {
		datetime := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		dayOfMonth := datetime.AddDate(0, 1, -1).Day()

		schedules, errBlock := p.prayTimeRepo.GetMonthlyPrayTime(ctx, city.Code, m)
		if errBlock != nil {
			err = errBlock
			return
		}

		for j := 0; j < dayOfMonth; j++ {
			dayNum := datetime.Day()
			schedule, ok := schedules[dayNum]
			if !ok {
				datetime = datetime.AddDate(0, 0, 1)
				continue
			}

			var day string
			switch lang {
			case "en":
				day = domain.DayEn[datetime.Weekday()]
			case "id":
				day = domain.DayId[datetime.Weekday()]
			case "ar":
				day = domain.DayAr[datetime.Weekday()]
			}
			result := domain.PrayTimeResponseDTO{
				Location: domain.CityResponseDTO{
					CityCode:  city.Code,
					City:      city.City,
					Latitude:  city.Latitude,
					Longitude: city.Longitude,
					Timezone:  city.Timezone,
				},
				Date: domain.DateResponseDTO{
					Day:  day,
					Date: datetime.Format("2006-01-02"),
				},
				Times: buildTimeResponseDTO(year, time.Month(m), dayNum, schedule, loc),
			}
			res = append(res, result)
			datetime = datetime.AddDate(0, 0, 1)
		}
	}

	return
}


func (p *prayTimeUsecase) GetDailyByCode(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}

	return p.getDaily(ctx, city, date, lang)
}

func (p *prayTimeUsecase) GetDailyByCodeStandard(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetDailyByCode(ctx, cityCode, date, lang)
	if err != nil {
		return
	}
	switch res.Location.Timezone {
	case 7:
		loc = time.FixedZone("UTC+7", 7*60*60)
	case 8:
		loc = time.FixedZone("UTC+8", 8*60*60)
	case 9:
		loc = time.FixedZone("UTC+9", 9*60*60)
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.In(loc).Format("15:04")
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.In(loc).Format("15:04")
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.In(loc).Format("15:04")
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.In(loc).Format("15:04")
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.In(loc).Format("15:04")
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.In(loc).Format("15:04")

	return
}

func (p *prayTimeUsecase) GetDailyByCodeISO8601(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetDailyByCode(ctx, cityCode, date, lang)
	if err != nil {
		return
	}
	switch res.Location.Timezone {
	case 7:
		loc = time.FixedZone("UTC+7", 7*60*60)
	case 8:
		loc = time.FixedZone("UTC+8", 8*60*60)
	case 9:
		loc = time.FixedZone("UTC+9", 9*60*60)
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)

	return
}

func (p *prayTimeUsecase) GetDailyByCityName(ctx context.Context, cityName string, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityName(ctx, cityName)
	if err != nil {
		return
	}

	return p.getDaily(ctx, city, date, lang)
}

func (p *prayTimeUsecase) GetDailyByCityNameStandard(ctx context.Context, cityName string, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetDailyByCityName(ctx, cityName, date, lang)
	if err != nil {
		return
	}
	switch res.Location.Timezone {
	case 7:
		loc = time.FixedZone("UTC+7", 7*60*60)
	case 8:
		loc = time.FixedZone("UTC+8", 8*60*60)
	case 9:
		loc = time.FixedZone("UTC+9", 9*60*60)
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.In(loc).Format("15:04")
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.In(loc).Format("15:04")
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.In(loc).Format("15:04")
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.In(loc).Format("15:04")
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.In(loc).Format("15:04")
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.In(loc).Format("15:04")

	return
}

func (p *prayTimeUsecase) GetDailyByCityNameISO8601(ctx context.Context, cityName string, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetDailyByCityName(ctx, cityName, date, lang)
	if err != nil {
		return
	}
	switch res.Location.Timezone {
	case 7:
		loc = time.FixedZone("UTC+7", 7*60*60)
	case 8:
		loc = time.FixedZone("UTC+8", 8*60*60)
	case 9:
		loc = time.FixedZone("UTC+9", 9*60*60)
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)

	return
}

func (p *prayTimeUsecase) GetMonthlyByCode(ctx context.Context, cityCode int, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}

	return p.getMonthly(ctx, city, year, month, lang)
}

func (p *prayTimeUsecase) GetMonthlyByCodeStandard(ctx context.Context, cityCode int, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetMonthlyByCode(ctx, cityCode, year, month, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format("15:04")
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format("15:04")
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format("15:04")
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format("15:04")
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format("15:04")
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format("15:04")
	}

	return
}

func (p *prayTimeUsecase) GetMonthlyByCodeISO8601(ctx context.Context, cityCode int, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetMonthlyByCode(ctx, cityCode, year, month, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)
	}

	return
}

func (p *prayTimeUsecase) GetMonthlyByCityName(ctx context.Context, cityName string, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityName(ctx, cityName)
	if err != nil {
		return
	}

	return p.getMonthly(ctx, city, year, month, lang)
}

func (p *prayTimeUsecase) GetMonthlyByCityNameStandard(ctx context.Context, cityName string, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetMonthlyByCityName(ctx, cityName, year, month, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format("15:04")
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format("15:04")
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format("15:04")
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format("15:04")
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format("15:04")
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format("15:04")
	}

	return
}

func (p *prayTimeUsecase) GetMonthlyByCityNameISO8601(ctx context.Context, cityName string, year int, month int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetMonthlyByCityName(ctx, cityName, year, month, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)
	}

	return
}

func (p *prayTimeUsecase) GetYearlyByCode(ctx context.Context, cityCode int, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}

	return p.getYearly(ctx, city, year, lang)
}

func (p *prayTimeUsecase) GetYearlyByCodeStandard(ctx context.Context, cityCode int, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetYearlyByCode(ctx, cityCode, year, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format("15:04")
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format("15:04")
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format("15:04")
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format("15:04")
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format("15:04")
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format("15:04")
	}

	return
}

func (p *prayTimeUsecase) GetYearlyByCodeISO8601(ctx context.Context, cityCode int, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetYearlyByCode(ctx, cityCode, year, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)
	}

	return
}

func (p *prayTimeUsecase) GetYearlyByCityName(ctx context.Context, cityName string, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	city, err := p.cityRepo.GetByCityName(ctx, cityName)
	if err != nil {
		return
	}

	return p.getYearly(ctx, city, year, lang)
}

func (p *prayTimeUsecase) GetYearlyByCityNameStandard(ctx context.Context, cityName string, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetYearlyByCityName(ctx, cityName, year, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format("15:04")
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format("15:04")
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format("15:04")
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format("15:04")
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format("15:04")
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format("15:04")
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format("15:04")
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format("15:04")
	}

	return
}

func (p *prayTimeUsecase) GetYearlyByCityNameISO8601(ctx context.Context, cityName string, year int, lang string) (res []domain.PrayTimeResponseDTO, err error) {
	var loc *time.Location

	res, err = p.GetYearlyByCityName(ctx, cityName, year, lang)
	if err != nil {
		return
	}

	for idx, _ := range res {
		switch res[idx].Location.Timezone {
		case 7:
			loc = time.FixedZone("UTC+7", 7*60*60)
		case 8:
			loc = time.FixedZone("UTC+8", 8*60*60)
		case 9:
			loc = time.FixedZone("UTC+9", 9*60*60)
		}
		timeImsak := time.Unix(int64(res[idx].Times.Imsak.(int)), 0)
		res[idx].Times.Imsak = timeImsak.In(loc).Format(time.RFC3339)
		timeFajr := time.Unix(int64(res[idx].Times.Fajr.(int)), 0)
		res[idx].Times.Fajr = timeFajr.In(loc).Format(time.RFC3339)
		timeSunrise := time.Unix(int64(res[idx].Times.Sunrise.(int)), 0)
		res[idx].Times.Sunrise = timeSunrise.In(loc).Format(time.RFC3339)
		timeDhuhr := time.Unix(int64(res[idx].Times.Dhuhr.(int)), 0)
		res[idx].Times.Dhuhr = timeDhuhr.In(loc).Format(time.RFC3339)
		timeAsr := time.Unix(int64(res[idx].Times.Asr.(int)), 0)
		res[idx].Times.Asr = timeAsr.In(loc).Format(time.RFC3339)
		timeMaghrib := time.Unix(int64(res[idx].Times.Maghrib.(int)), 0)
		res[idx].Times.Maghrib = timeMaghrib.In(loc).Format(time.RFC3339)
		timeIsha := time.Unix(int64(res[idx].Times.Isha.(int)), 0)
		res[idx].Times.Isha = timeIsha.In(loc).Format(time.RFC3339)
		timeMidnight := time.Unix(int64(res[idx].Times.Midnight.(int)), 0)
		res[idx].Times.Midnight = timeMidnight.In(loc).Format(time.RFC3339)
	}

	return
}

func (p *prayTimeUsecase) CreatePrayTime(ctx context.Context, prayTime *domain.PrayTime) (err error) {
	datetime, err := time.Parse("2006-01-02", prayTime.Date)
	if err != nil {
		datetime, err = time.Parse("20060102", prayTime.Date)
		if err != nil {
			return err
		}
	}
	city, err := p.cityRepo.GetByCityCode(ctx, prayTime.CityCode)
	if err != nil {
		return err
	}
	loc := getTimeLocation(city.Timezone)

	schedule := domain.PrayTimeSchedule{
		Imsak:    time.Unix(int64(prayTime.Imsak), 0).In(loc).Format("15:04"),
		Fajr:     time.Unix(int64(prayTime.Fajr), 0).In(loc).Format("15:04"),
		Sunrise:  time.Unix(int64(prayTime.Sunrise), 0).In(loc).Format("15:04"),
		Dhuhr:    time.Unix(int64(prayTime.Dhuhr), 0).In(loc).Format("15:04"),
		Asr:      time.Unix(int64(prayTime.Asr), 0).In(loc).Format("15:04"),
		Maghrib:  time.Unix(int64(prayTime.Maghrib), 0).In(loc).Format("15:04"),
		Isha:     time.Unix(int64(prayTime.Isha), 0).In(loc).Format("15:04"),
		Midnight: time.Unix(int64(prayTime.Midnight), 0).In(loc).Format("15:04"),
	}

	return p.prayTimeRepo.SaveDaily(ctx, prayTime.CityCode, int(datetime.Month()), datetime.Day(), &schedule)
}

func (p *prayTimeUsecase) CreatePrayTimeMonthlyAPIMyQuran(ctx context.Context, cityCode int, year int, month int) (err error) {
	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}

	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.Header.SetContentType("application/json")
	req.SetRequestURI("https://api.myquran.com/v1/sholat/jadwal/" + strconv.FormatInt(int64(city.IDmyQuran), 10) + "/" + strconv.FormatInt(int64(year), 10) + "/" + fmt.Sprintf("%02d", month))
	if err = a.Parse(); err != nil {
		return
	}

	code, body, errs := a.Bytes()
	log.Debug(errs)

	var response apiMyQuranMonthlyResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	if code == fiber.StatusOK {
		schedules := make(map[int]domain.PrayTimeSchedule, len(response.Data.Jadwal))
		for _, pray := range response.Data.Jadwal {
			date, errDate := time.Parse("2006-01-02", pray.Date)
			if errDate != nil {
				log.Error(errDate)
				continue
			}

			sunriseTime, errSunrise := time.Parse("15:04", pray.Terbit)
			if errSunrise != nil {
				log.Error(errSunrise)
				continue
			}
			maghribTime, errMaghrib := time.Parse("15:04", pray.Maghrib)
			if errMaghrib != nil {
				log.Error(errMaghrib)
				continue
			}

			nextSunriseTime := sunriseTime.AddDate(0, 0, 1)
			diff := nextSunriseTime.Sub(maghribTime).Seconds() / 2
			midnightTime := maghribTime.Add(time.Second * time.Duration(diff))

			schedule := domain.PrayTimeSchedule{
				Imsak:    pray.Imsak,
				Fajr:     pray.Subuh,
				Sunrise:  pray.Terbit,
				Dhuhr:    pray.Dzuhur,
				Asr:      pray.Ashar,
				Maghrib:  pray.Maghrib,
				Isha:     pray.Isya,
				Midnight: midnightTime.Format("15:04"),
			}
			schedules[date.Day()] = schedule
		}

		err = p.prayTimeRepo.SaveMonthly(ctx, cityCode, month, schedules)
		if err != nil {
			log.Error(err)
		}
	} else {
		err = errors.New("response " + strconv.FormatInt(int64(code), 10))
	}

	return
}

