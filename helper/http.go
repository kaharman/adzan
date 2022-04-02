package helper

import (
	"github.com/gofiber/fiber/v2"
)

type errorStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HttpSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

func HttpError(c *fiber.Ctx, httpStatus int, message string, data interface{}) error {
	return c.Status(httpStatus).JSON(&errorStruct{
		Message: message,
		Data:    data,
	})
}
