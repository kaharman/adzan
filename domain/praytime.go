package domain

import "context"

type PrayTime struct {
	CityCode int    `db:"city_code"`
	Date     string `db:"date"`
	Imsak    int    `db:"imsak"`
	Fajr     int    `db:"fajr"`
	Sunrise  int    `db:"sunrise"`
	Dhuhr    int    `db:"dhuhr"`
	Asr      int    `db:"asr"`
	Maghrib  int    `db:"maghrib"`
	Isha     int    `db:"isha"`
	Midnight int    `db:"midnight"`
}

type DateResponseDTO struct {
	Day  string `json:"day"`
	Date string `json:"date"`
}

type TimeResponseDTO struct {
	Imsak    interface{} `json:"imsak"`
	Fajr     interface{} `json:"fajr"`
	Sunrise  interface{} `json:"sunrise"`
	Dhuhr    interface{} `json:"dhuhr"`
	Asr      interface{} `json:"asr"`
	Maghrib  interface{} `json:"maghrib"`
	Isha     interface{} `json:"isha"`
	Midnight interface{} `json:"midnight"`
}

type PrayTimeResponseDTO struct {
	Location CityResponseDTO `json:"location"`
	Date     DateResponseDTO `json:"date"`
	Times    TimeResponseDTO `json:"times"`
}

type PrayTimeSchedule struct {
	Imsak    string `json:"imsak"`
	Fajr     string `json:"fajr"`
	Sunrise  string `json:"sunrise"`
	Dhuhr    string `json:"dhuhr"`
	Asr      string `json:"asr"`
	Maghrib  string `json:"maghrib"`
	Isha     string `json:"isha"`
	Midnight string `json:"midnight"`
}

type PrayTimeUsecase interface {
	GetDailyByCode(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetDailyByCodeStandard(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetDailyByCodeISO8601(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetDailyByCityName(ctx context.Context, cityName string, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetDailyByCityNameStandard(ctx context.Context, cityName string, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetDailyByCityNameISO8601(ctx context.Context, cityName string, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetMonthlyByCode(ctx context.Context, cityCode int, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetMonthlyByCodeStandard(ctx context.Context, cityCode int, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetMonthlyByCodeISO8601(ctx context.Context, cityCode int, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetMonthlyByCityName(ctx context.Context, cityName string, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetMonthlyByCityNameStandard(ctx context.Context, cityName string, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetMonthlyByCityNameISO8601(ctx context.Context, cityName string, year int, month int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCode(ctx context.Context, cityCode int, year int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCodeStandard(ctx context.Context, cityCode int, year int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCodeISO8601(ctx context.Context, cityCode int, year int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCityName(ctx context.Context, cityName string, year int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCityNameStandard(ctx context.Context, cityName string, year int, lang string) (res []PrayTimeResponseDTO, err error)
	GetYearlyByCityNameISO8601(ctx context.Context, cityName string, year int, lang string) (res []PrayTimeResponseDTO, err error)
	CreatePrayTime(ctx context.Context, pray *PrayTime) (err error)
	CreatePrayTimeMonthlyAPIMyQuran(ctx context.Context, cityCode int, year int, month int) (err error)
}

type PrayTimeRepository interface {
	GetPrayTime(ctx context.Context, cityCode int, month int, day int) (res PrayTimeSchedule, err error)
	GetMonthlyPrayTime(ctx context.Context, cityCode int, month int) (res map[int]PrayTimeSchedule, err error)
	SaveDaily(ctx context.Context, cityCode int, month int, day int, schedule *PrayTimeSchedule) (err error)
	SaveMonthly(ctx context.Context, cityCode int, month int, schedules map[int]PrayTimeSchedule) (err error)
	Delete(ctx context.Context, cityCode int, month int, day int) (err error)
}


