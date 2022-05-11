package domain

import (
	"context"
	"time"
)

type Hijri struct {
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	MonthName string `json:"month_name"`
	String    string `json:"string"`
	Notes     string `json:"notes"`
}

type HijriUsecase interface {
	GregorianToHijri(ctx context.Context, date time.Time, lang string) (res Hijri, err error)
}

type HijriRepository interface {
	GregorianToHijri(ctx context.Context, date time.Time) (res Hijri, err error)
}
