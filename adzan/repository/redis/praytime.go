package redis

import (
	"adzan/domain"
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type prayTime struct {
	Imsak    int `redis:"imsak"`
	Fajr     int `redis:"fajr"`
	Sunrise  int `redis:"sunrise"`
	Dhuhr    int `redis:"dhuhr"`
	Asr      int `redis:"asr"`
	Maghrib  int `redis:"maghrib"`
	Isha     int `redis:"isha"`
	Midnight int `redis:"midnight"`
}

type redisPrayTimeRepository struct {
	Conn *redis.Client
}

func NewRedisPrayTimeRepository(Conn *redis.Client) domain.PrayTimeRepository {
	return &redisPrayTimeRepository{Conn}
}

func (r *redisPrayTimeRepository) GetPrayTime(ctx context.Context, cityCode int, date string) (res domain.PrayTime, err error) {
	var prayTime prayTime
	err = r.Conn.HGetAll(ctx, "praytime:"+strconv.FormatInt(int64(cityCode), 10)+":"+date).Scan(&prayTime)
	if err != nil {
		return
	}
	res = domain.PrayTime{
		CityCode: cityCode,
		Date:     date,
		Imsak:    prayTime.Imsak,
		Fajr:     prayTime.Fajr,
		Sunrise:  prayTime.Sunrise,
		Dhuhr:    prayTime.Dhuhr,
		Asr:      prayTime.Asr,
		Maghrib:  prayTime.Maghrib,
		Isha:     prayTime.Isha,
		Midnight: prayTime.Midnight,
	}

	return
}

func (r *redisPrayTimeRepository) Create(ctx context.Context, pray *domain.PrayTime) (err error) {
	_, err = r.Conn.HSet(ctx, "praytime:"+strconv.FormatInt(int64(pray.CityCode), 10)+":"+pray.Date, "imsak", pray.Imsak, "fajr", pray.Fajr, "sunrise", pray.Sunrise, "dhuhr", pray.Dhuhr, "asr", pray.Asr, "maghrib", pray.Maghrib, "isha", pray.Isha, "midnight", pray.Midnight).Result()
	if err != nil {
		return
	}

	return
}

func (r *redisPrayTimeRepository) Update(ctx context.Context, prayTime *domain.PrayTime) (err error) {
	return r.Create(ctx, prayTime)
}

func (r *redisPrayTimeRepository) Delete(ctx context.Context, cityCode int, date string) (err error) {
	_, err = r.Conn.Del(ctx, "praytime:"+strconv.FormatInt(int64(cityCode), 10)+":"+date).Result()
	return
}
