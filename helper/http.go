package helper

import (
	"log/slog"

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
	if httpStatus >= 500 {
		slog.Error("HTTP error response",
			"status", httpStatus,
			"message", message,
			"method", c.Method(),
			"path", c.OriginalURL(),
			"ip", c.IP(),
			"data", data,
		)
	} else if httpStatus >= 400 {
		slog.Warn("HTTP warning response",
			"status", httpStatus,
			"message", message,
			"method", c.Method(),
			"path", c.OriginalURL(),
			"ip", c.IP(),
			"data", data,
		)
	}

	return c.Status(httpStatus).JSON(&errorStruct{
		Message: message,
		Data:    data,
	})
}
