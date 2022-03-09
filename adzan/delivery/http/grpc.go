package http

import (
	"embed"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/utils"
)

//go:embed grpcdoc/*
var gRPC embed.FS

func RouterGRPC(app *fiber.App) {
	app.Use("/grpcdoc", func(c *fiber.Ctx) error {
		originalUrl := utils.ImmutableString(c.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalUrl)

		if !strings.HasSuffix(originalUrl, "/") && !extMatch {
			return c.Redirect(originalUrl + "/")
		}
		return c.Next()
	}, filesystem.New(filesystem.Config{
		Root:       http.FS(gRPC),
		PathPrefix: "grpcdoc",
	}))
}
