package http

import (
	"adzan/domain"
	"adzan/helper"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type AdzanHandler struct {
	PrayTimeUsecase domain.PrayTimeUsecase
	CityUsecase     domain.CityUsecase
}

func RouterAPI(app *fiber.App, prayTime domain.PrayTimeUsecase, city domain.CityUsecase) {
	// middleware := &MissIndonesiaMiddleware{}
	handler := &AdzanHandler{
		PrayTimeUsecase: prayTime,
		CityUsecase:     city,
	}

	app.Get("/today", handler.GetTodayPrayTime)
	app.Get("/daily", handler.GetPrayTimeDaily)
	app.Get("/monthly", handler.GetPrayTimeMonthly)
	app.Get("/yearly", handler.GetPrayTimeYearly)
	app.Post("/apimyquranmonthly", handler.CreatePrayTimeMonthlyAPIMyQuran)

	app.Get("/city", handler.GetCity)
}

func (adzan *AdzanHandler) GetTodayPrayTime(c *fiber.Ctx) (err error) {
	var cityCode int64

	format := c.Query("format", "standard")
	if format != "standard" && format != "unixtime" && format != "iso8601" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format", nil)
	}
	language := c.Query("language", "en")
	if language != "en" && language != "id" && language != "ar" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid language", nil)
	}
	cityCodeStr := c.Query("city_code", "")
	cityName := c.Query("city", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
		}
		today := time.Now().Format("2006-01-02")
		var response domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetDailyByCodeStandard(c.Context(), int(cityCode), today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetDailyByCode(c.Context(), int(cityCode), today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetDailyByCodeISO8601(c.Context(), int(cityCode), today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	} else if cityName != "" {
		today := time.Now().Format("2006-01-02")
		var response domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityNameStandard(c.Context(), cityName, today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityName(c.Context(), cityName, today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityNameISO8601(c.Context(), cityName, today, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	}

	return helper.HttpError(c, fiber.StatusBadRequest, "Parameters required", nil)
}

func (adzan *AdzanHandler) GetPrayTimeDaily(c *fiber.Ctx) (err error) {
	var cityCode int64

	dateStr := c.Query("date")
	if dateStr == "" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid date", nil)
	}
	_, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid date", nil)
	}
	format := c.Query("format", "standard")
	if format != "standard" && format != "unixtime" && format != "iso8601" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format", nil)
	}
	language := c.Query("language", "en")
	if language != "en" && language != "id" && language != "ar" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid language", nil)
	}
	cityCodeStr := c.Query("city_code", "")
	cityName := c.Query("city", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
		}
		var response domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetDailyByCodeStandard(c.Context(), int(cityCode), dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetDailyByCode(c.Context(), int(cityCode), dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetDailyByCodeISO8601(c.Context(), int(cityCode), dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	} else if cityName != "" {
		var response domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityNameStandard(c.Context(), cityName, dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityName(c.Context(), cityName, dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetDailyByCityNameISO8601(c.Context(), cityName, dateStr, language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	}

	return helper.HttpError(c, fiber.StatusBadRequest, "Parameters required", nil)
}

func (adzan *AdzanHandler) GetPrayTimeMonthly(c *fiber.Ctx) (err error) {
	var cityCode int64

	yearStr := c.Query("year")
	if yearStr == "" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	if year < 2022 {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	monthStr := c.Query("month")
	if monthStr == "" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid month", nil)
	}
	month, err := strconv.ParseInt(monthStr, 10, 64)
	if err != nil {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid month", nil)
	}
	if month < 1 || month > 12 {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid month", nil)
	}
	format := c.Query("format", "standard")
	if format != "standard" && format != "unixtime" && format != "iso8601" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format", nil)
	}
	language := c.Query("language", "en")
	if language != "en" && language != "id" && language != "ar" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid language", nil)
	}
	cityCodeStr := c.Query("city_code", "")
	cityName := c.Query("city", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
		}
		var response []domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCodeStandard(c.Context(), int(cityCode), int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCode(c.Context(), int(cityCode), int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCodeISO8601(c.Context(), int(cityCode), int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	} else if cityName != "" {
		var response []domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCityNameStandard(c.Context(), cityName, int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCityName(c.Context(), cityName, int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetMonthlyByCityNameISO8601(c.Context(), cityName, int(year), int(month), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	}

	return helper.HttpError(c, fiber.StatusBadRequest, "Parameters required", nil)
}

func (adzan *AdzanHandler) GetPrayTimeYearly(c *fiber.Ctx) (err error) {
	var cityCode int64

	yearStr := c.Query("year")
	if yearStr == "" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	if year < 2022 {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid year", nil)
	}
	format := c.Query("format", "standard")
	if format != "standard" && format != "unixtime" && format != "iso8601" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format", nil)
	}
	language := c.Query("language", "en")
	if language != "en" && language != "id" && language != "ar" {
		return helper.HttpError(c, fiber.StatusBadRequest, "Invalid language", nil)
	}
	cityCodeStr := c.Query("city_code", "")
	cityName := c.Query("city", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
		}
		var response []domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCodeStandard(c.Context(), int(cityCode), int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCode(c.Context(), int(cityCode), int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCodeISO8601(c.Context(), int(cityCode), int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city_code", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	} else if cityName != "" {
		var response []domain.PrayTimeResponseDTO

		switch format {
		case "standard":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCityNameStandard(c.Context(), cityName, int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "unixtime":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCityName(c.Context(), cityName, int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		case "iso8601":
			response, err = adzan.PrayTimeUsecase.GetYearlyByCityNameISO8601(c.Context(), cityName, int(year), language)
			if err != nil {
				if err.Error() == "redis: nil" {
					return helper.HttpError(c, fiber.StatusBadRequest, "Invalid city", nil)
				} else {
					log.Debug(err)
					return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
				}
			}
		}
		return helper.HttpSuccess(c, response)
	}

	return helper.HttpError(c, fiber.StatusBadRequest, "Parameters required", nil)
}

func (adzan *AdzanHandler) CreatePrayTimeMonthlyAPIMyQuran(c *fiber.Ctx) (err error) {
	var cityCode, year, month int64

	yearStr := c.FormValue("year", "")
	if yearStr != "" {
		year, err = strconv.ParseInt(yearStr, 10, 64)
		if err != nil || year < 2022 {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format year", nil)
		}
	} else {
		return helper.HttpError(c, fiber.StatusBadRequest, "Parameter year required", nil)
	}
	monthStr := c.FormValue("month", "")
	if monthStr != "" {
		month, err = strconv.ParseInt(monthStr, 10, 64)
		if err != nil || month < 1 || month > 12 {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format month", nil)
		}
	} else {
		return helper.HttpError(c, fiber.StatusBadRequest, "Parameter month required", nil)
	}
	cityCodeStr := c.FormValue("city_code", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format city_code", nil)
		}
	} else {
		return helper.HttpError(c, fiber.StatusBadRequest, "Parameter coty_code required", nil)
	}

	return adzan.PrayTimeUsecase.CreatePrayTimeMonthlyAPIMyQuran(c.Context(), int(cityCode), int(year), int(month))
}

func (adzan *AdzanHandler) GetCity(c *fiber.Ctx) (err error) {
	var latitude, longitude float64
	var cityCode int64

	latitudeStr := c.Query("latitude", "")
	longitudeStr := c.Query("longitude", "")
	cityCodeStr := c.Query("city_code", "")
	if cityCodeStr != "" {
		cityCode, err = strconv.ParseInt(cityCodeStr, 10, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format city_code", nil)
		}
		response, err := adzan.CityUsecase.GetByCityCode(c.Context(), int(cityCode))
		if err != nil {
			return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
		}
		return helper.HttpSuccess(c, response)
	} else if latitudeStr != "" && longitudeStr != "" {
		latitude, err = strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format latitude", nil)
		}
		longitude, err = strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			return helper.HttpError(c, fiber.StatusBadRequest, "Invalid format longitude", nil)
		}
		response, err := adzan.CityUsecase.GetByLatLong(c.Context(), float32(latitude), float32(longitude))
		if err != nil {
			return helper.HttpError(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
		}
		return helper.HttpSuccess(c, response)
	}

	return helper.HttpError(c, fiber.StatusBadRequest, "Parameters required", nil)
}
