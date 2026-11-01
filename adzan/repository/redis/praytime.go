package redis

import (
	"adzan/domain"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type redisPrayTimeRepository struct {
	Conn *redis.Client
}

func NewRedisPrayTimeRepository(Conn *redis.Client) domain.PrayTimeRepository {
	return &redisPrayTimeRepository{Conn}
}

func keyMonthly(cityCode int, month int) string {
	return fmt.Sprintf("praytime:%d:%02d", cityCode, month)
}

func (r *redisPrayTimeRepository) GetPrayTime(ctx context.Context, cityCode int, month int, day int) (res domain.PrayTimeSchedule, err error) {
	key := keyMonthly(cityCode, month)
	field := fmt.Sprintf("%02d", day)

	val, err := r.Conn.HGet(ctx, key, field).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(val), &res)
	return
}

func (r *redisPrayTimeRepository) GetMonthlyPrayTime(ctx context.Context, cityCode int, month int) (res map[int]domain.PrayTimeSchedule, err error) {
	key := keyMonthly(cityCode, month)

	vals, err := r.Conn.HGetAll(ctx, key).Result()
	if err != nil {
		return
	}
	if len(vals) == 0 {
		return nil, redis.Nil
	}

	res = make(map[int]domain.PrayTimeSchedule, len(vals))
	for field, jsonVal := range vals {
		dayNum, errParse := strconv.Atoi(field)
		if errParse != nil {
			continue
		}
		var schedule domain.PrayTimeSchedule
		if errJson := json.Unmarshal([]byte(jsonVal), &schedule); errJson == nil {
			res[dayNum] = schedule
		}
	}

	return
}

func (r *redisPrayTimeRepository) SaveDaily(ctx context.Context, cityCode int, month int, day int, schedule *domain.PrayTimeSchedule) (err error) {
	key := keyMonthly(cityCode, month)
	field := fmt.Sprintf("%02d", day)

	data, err := json.Marshal(schedule)
	if err != nil {
		return err
	}

	return r.Conn.HSet(ctx, key, field, string(data)).Err()
}

func (r *redisPrayTimeRepository) SaveMonthly(ctx context.Context, cityCode int, month int, schedules map[int]domain.PrayTimeSchedule) (err error) {
	key := keyMonthly(cityCode, month)

	values := make(map[string]interface{}, len(schedules))
	for day, schedule := range schedules {
		data, errMarshal := json.Marshal(schedule)
		if errMarshal != nil {
			return errMarshal
		}
		values[fmt.Sprintf("%02d", day)] = string(data)
	}

	if len(values) == 0 {
		return nil
	}

	return r.Conn.HSet(ctx, key, values).Err()
}

func (r *redisPrayTimeRepository) Delete(ctx context.Context, cityCode int, month int, day int) (err error) {
	key := keyMonthly(cityCode, month)
	field := fmt.Sprintf("%02d", day)
	return r.Conn.HDel(ctx, key, field).Err()
}

