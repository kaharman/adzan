package lib

import (
	"adzan/domain"
	"context"
	"time"

	"github.com/hablullah/go-hijri"
)

type libHijriRepository struct{}

func NewLibHijriRepository() domain.HijriRepository {
	return &libHijriRepository{}
}

func (r *libHijriRepository) GregorianToHijri(ctx context.Context, date time.Time) (res domain.Hijri, err error) {
	ummAlQuraDate, _ := hijri.CreateUmmAlQuraDate(date)

	res = domain.Hijri{
		Day:   int(ummAlQuraDate.Day),
		Month: int(ummAlQuraDate.Month),
		Year:  int(ummAlQuraDate.Year),
		Notes: "Umm al-Qura",
	}

	return
}
