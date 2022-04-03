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

type PrayTimeUsecase interface {
	GetTodayByCode(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetTodayByCodeStandard(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	GetTodayByCodeISO8601(ctx context.Context, cityCode int, date string, lang string) (res PrayTimeResponseDTO, err error)
	CreatePrayTime(ctx context.Context, pray *PrayTime) (err error)
	CreatePrayTimeMonthlyAPIMyQuran(ctx context.Context, cityCode int, year int, month int) (err error)
}

type PrayTimeRepository interface {
	GetPrayTime(ctx context.Context, cityCode int, date string) (res PrayTime, err error)
	Create(ctx context.Context, prayTime *PrayTime) (err error)
	Update(ctx context.Context, prayTime *PrayTime) (err error)
	Delete(ctx context.Context, cityCode int, date string) (err error)
}
