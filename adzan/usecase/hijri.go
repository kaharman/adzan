package usecase

import (
	"adzan/domain"
	"context"
	"fmt"
	"strconv"
	"time"
)

type hijriUsecase struct {
	hijriRepo domain.HijriRepository
}

func NewHijriUsecase(hijriRepo domain.HijriRepository) domain.HijriUsecase {
	return &hijriUsecase{
		hijriRepo: hijriRepo,
	}
}

func (c *hijriUsecase) GregorianToHijri(ctx context.Context, date time.Time, lang string) (res domain.Hijri, err error) {
	res, err = c.hijriRepo.GregorianToHijri(ctx, date)
	if err != nil {
		return
	}

	switch lang {
	case "en":
		res.MonthName = domain.HijriMonthEn[res.Month-1]
		res.String = fmt.Sprintf("%d %s %d H", res.Day, domain.HijriMonthEn[res.Month-1], res.Year)
	case "id":
		res.MonthName = domain.HijriMonthId[res.Month-1]
		res.String = fmt.Sprintf("%d %s %d H", res.Day, domain.HijriMonthId[res.Month-1], res.Year)
	case "ar":
		dayStr := strconv.FormatInt(int64(res.Day), 10)
		yearStr := strconv.FormatInt(int64(res.Year), 10)
		res.MonthName = domain.HijriMonthAr[res.Month-1]
		res.String = fmt.Sprintf("%s %s %s \ufeeb", domain.ConvertDigitEnglishToArabic(dayStr), domain.HijriMonthAr[res.Month-1], domain.ConvertDigitEnglishToArabic(yearStr))
	}

	return
}
