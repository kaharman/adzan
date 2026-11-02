package http

import (
	"embed"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/utils"
	"github.com/spf13/viper"
)

//go:embed openapi/*
var openAPI embed.FS

func RouterOpenAPI(router fiber.Router) {
	// Serve openapi.yaml dynamically following server.base_path
	router.Get("/openapi/openapi.yaml", func(c *fiber.Ctx) error {
		basePath := strings.TrimSpace(viper.GetString("server.base_path"))
		if basePath != "" && !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		basePath = strings.TrimRight(basePath, "/")
		serverURL := basePath
		if serverURL == "" {
			serverURL = "/"
		}

		content, err := openAPI.ReadFile("openapi/openapi.yaml")
		if err != nil {
			return err
		}

		replaced := strings.Replace(
			string(content),
			"servers:\n  - url: /",
			fmt.Sprintf("servers:\n  - url: %s", serverURL),
			1,
		)

		c.Set("Content-Type", "application/yaml")
		return c.SendString(replaced)
	})

	router.Use("/openapi", func(c *fiber.Ctx) error {
		originalUrl := utils.ImmutableString(c.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalUrl)

		if !strings.HasSuffix(originalUrl, "/") && !extMatch {
			return c.Redirect(originalUrl + "/")
		}
		return c.Next()
	}, filesystem.New(filesystem.Config{
		Root:       http.FS(openAPI),
		PathPrefix: "openapi",
	}))
}

