package http

import (
	"adzan/domain"
	"adzan/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AdzanHandler struct {
	CityUsecase domain.CityUsecase
}

func RouterAPI(app *fiber.App, city domain.CityUsecase) {
	// middleware := &MissIndonesiaMiddleware{}
	handler := &AdzanHandler{
		CityUsecase: city,
	}

	app.Get("/today", handler.GetTodayPrayTime)
	app.Get("/daily", handler.GetPrayTimeDaily)
	app.Get("/monthly", handler.GetPrayTimeMonthly)
	app.Get("/yearly", handler.GetPrayTimeYearly)

	app.Get("/city", handler.GetCity)
}

func (adzan *AdzanHandler) GetTodayPrayTime(c *fiber.Ctx) (err error) {
	return
}

func (adzan *AdzanHandler) GetPrayTimeDaily(c *fiber.Ctx) (err error) {
	return
}

func (adzan *AdzanHandler) GetPrayTimeMonthly(c *fiber.Ctx) (err error) {
	return
}

func (adzan *AdzanHandler) GetPrayTimeYearly(c *fiber.Ctx) (err error) {
	return
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
