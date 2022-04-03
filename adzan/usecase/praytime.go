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

func (p *prayTimeUsecase) GetTodayByCode(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	datetime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return
	}

	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}

	prayTime, err := p.prayTimeRepo.GetPrayTime(ctx, cityCode, datetime.Format("20060102"))
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
		Times: domain.TimeResponseDTO{
			Imsak:    prayTime.Imsak,
			Fajr:     prayTime.Fajr,
			Sunrise:  prayTime.Sunrise,
			Dhuhr:    prayTime.Dhuhr,
			Asr:      prayTime.Asr,
			Maghrib:  prayTime.Maghrib,
			Isha:     prayTime.Isha,
			Midnight: prayTime.Midnight,
		},
	}

	return
}

func (p *prayTimeUsecase) GetTodayByCodeStandard(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	res, err = p.GetTodayByCode(ctx, cityCode, date, lang)
	if err != nil {
		return
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.Format("15:04")
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.Format("15:04")
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.Format("15:04")
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.Format("15:04")
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.Format("15:04")
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.Format("15:04")
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.Format("15:04")
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.Format("15:04")

	return
}

func (p *prayTimeUsecase) GetTodayByCodeISO8601(ctx context.Context, cityCode int, date string, lang string) (res domain.PrayTimeResponseDTO, err error) {
	res, err = p.GetTodayByCode(ctx, cityCode, date, lang)
	if err != nil {
		return
	}
	timeImsak := time.Unix(int64(res.Times.Imsak.(int)), 0)
	res.Times.Imsak = timeImsak.Format(time.RFC3339)
	timeFajr := time.Unix(int64(res.Times.Fajr.(int)), 0)
	res.Times.Fajr = timeFajr.Format(time.RFC3339)
	timeSunrise := time.Unix(int64(res.Times.Sunrise.(int)), 0)
	res.Times.Sunrise = timeSunrise.Format(time.RFC3339)
	timeDhuhr := time.Unix(int64(res.Times.Dhuhr.(int)), 0)
	res.Times.Dhuhr = timeDhuhr.Format(time.RFC3339)
	timeAsr := time.Unix(int64(res.Times.Asr.(int)), 0)
	res.Times.Asr = timeAsr.Format(time.RFC3339)
	timeMaghrib := time.Unix(int64(res.Times.Maghrib.(int)), 0)
	res.Times.Maghrib = timeMaghrib.Format(time.RFC3339)
	timeIsha := time.Unix(int64(res.Times.Isha.(int)), 0)
	res.Times.Isha = timeIsha.Format(time.RFC3339)
	timeMidnight := time.Unix(int64(res.Times.Midnight.(int)), 0)
	res.Times.Midnight = timeMidnight.Format(time.RFC3339)

	return
}

func (p *prayTimeUsecase) CreatePrayTime(ctx context.Context, prayTime *domain.PrayTime) (err error) {
	return p.prayTimeRepo.Create(ctx, prayTime)
}

func (p *prayTimeUsecase) CreatePrayTimeMonthlyAPIMyQuran(ctx context.Context, cityCode int, year int, month int) (err error) {
	var loc *time.Location

	city, err := p.cityRepo.GetByCityCode(ctx, cityCode)
	if err != nil {
		return
	}
	switch city.Timezone {
	case 7:
		loc = time.FixedZone("UTC+7", 7*60*60)
	case 8:
		loc = time.FixedZone("UTC+8", 8*60*60)
	case 9:
		loc = time.FixedZone("UTC+9", 9*60*60)
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
	// log.Debug(code)
	// log.Debug(string(body))
	log.Debug(errs)

	var response apiMyQuranMonthlyResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	if code == fiber.StatusOK {
		// log.Debug(response)
		for _, pray := range response.Data.Jadwal {
			date, err := time.Parse("2006-01-02", pray.Date)
			if err != nil {
				log.Error(err)
				continue
			}
			imsakTime, err := time.Parse("15:04", pray.Imsak)
			if err != nil {
				log.Error(err)
				continue
			}
			imsak := time.Date(date.Year(), date.Month(), date.Day(), imsakTime.Hour(), imsakTime.Minute(), 0, 0, loc)
			fajrTime, err := time.Parse("15:04", pray.Subuh)
			if err != nil {
				log.Error(err)
				continue
			}
			fajr := time.Date(date.Year(), date.Month(), date.Day(), fajrTime.Hour(), fajrTime.Minute(), 0, 0, loc)
			sunriseTime, err := time.Parse("15:04", pray.Terbit)
			if err != nil {
				log.Error(err)
				continue
			}
			sunrise := time.Date(date.Year(), date.Month(), date.Day(), sunriseTime.Hour(), sunriseTime.Minute(), 0, 0, loc)
			dhuhrTime, err := time.Parse("15:04", pray.Dzuhur)
			if err != nil {
				log.Error(err)
				continue
			}
			dhuhr := time.Date(date.Year(), date.Month(), date.Day(), dhuhrTime.Hour(), dhuhrTime.Minute(), 0, 0, loc)
			asrTime, err := time.Parse("15:04", pray.Ashar)
			if err != nil {
				log.Error(err)
				continue
			}
			asr := time.Date(date.Year(), date.Month(), date.Day(), asrTime.Hour(), asrTime.Minute(), 0, 0, loc)
			maghribTime, err := time.Parse("15:04", pray.Maghrib)
			if err != nil {
				log.Error(err)
				continue
			}
			maghrib := time.Date(date.Year(), date.Month(), date.Day(), maghribTime.Hour(), maghribTime.Minute(), 0, 0, loc)
			ishaTime, err := time.Parse("15:04", pray.Isya)
			if err != nil {
				log.Error(err)
				continue
			}
			isha := time.Date(date.Year(), date.Month(), date.Day(), ishaTime.Hour(), ishaTime.Minute(), 0, 0, loc)
			nextSunriseTime := sunriseTime.AddDate(0, 0, 1)
			diff := nextSunriseTime.Sub(maghribTime).Seconds() / 2
			prayTime := domain.PrayTime{
				CityCode: cityCode,
				Date:     date.Format("20060102"),
				Imsak:    int(imsak.Unix()),
				Fajr:     int(fajr.Unix()),
				Sunrise:  int(sunrise.Unix()),
				Dhuhr:    int(dhuhr.Unix()),
				Asr:      int(asr.Unix()),
				Maghrib:  int(maghrib.Unix()),
				Isha:     int(isha.Unix()),
				Midnight: int(maghrib.Add(time.Second * time.Duration(diff)).Unix()),
			}
			log.Debug(prayTime)

			err = p.prayTimeRepo.Create(ctx, &prayTime)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	} else {
		err = errors.New("response " + strconv.FormatInt(int64(code), 10))
	}

	return
}
