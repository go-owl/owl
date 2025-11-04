package main

import (
	"log"
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

func main() {
	// Example 1: Default body limit (10MB)
	app1 := owl.New(owl.AppConfig{
		Name:    "DefaultLimitAPI",
		Version: "1.0.0",
		// BodyLimit not set = 10MB default
	})

	// Example 2: Custom body limit (1MB)
	app2 := owl.New(owl.AppConfig{
		Name:      "SmallLimitAPI",
		Version:   "1.0.0",
		BodyLimit: 1 * owl.MB, // 1MB - easy to read!
	})

	// Example 3: Large body limit (50MB) for file uploads
	app3 := owl.New(owl.AppConfig{
		Name:      "LargeLimitAPI",
		Version:   "1.0.0",
		BodyLimit: 50 * owl.MB, // 50MB
	})

	// Example 4: No body limit (unlimited)
	app4 := owl.New(owl.AppConfig{
		Name:      "UnlimitedAPI",
		Version:   "1.0.0",
		BodyLimit: 0, // 0 = unlimited (not recommended for production)
	})

	// Use the custom limit example
	app := app2

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	api := app.Group("/api")

	// Upload endpoint - will be limited by body size
	api.POST("/upload", func(c *owl.Ctx) error {
		var data map[string]interface{}

		// This will fail if request body > 1MB
		if err := c.Bind().JSON(&data); err != nil {
			return owl.NewHTTPError(http.StatusBadRequest, "Request body too large or invalid JSON")
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "Data uploaded successfully",
			"size":    len(c.Request.Header.Get("Content-Length")),
		})
	})

	// Health check
	api.GET("/health", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status":     "ok",
			"body_limit": "1MB",
		})
	})

	log.Println("Body limit examples:")
	log.Println("- Default:   10MB (app1)")
	log.Println("- Custom:    1MB  (app2) <- currently running")
	log.Println("- Large:     50MB (app3)")
	log.Println("- Unlimited: 0    (app4)")
	log.Println()
	log.Println("Try uploading data:")
	log.Println("  curl -X POST http://localhost:8080/api/upload \\")
	log.Println("    -H 'Content-Type: application/json' \\")
	log.Println("    -d '{\"data\":\"your data here\"}'")
	log.Println()

	// Prevent unused variable warnings
	_, _, _, _ = app1, app2, app3, app4

	log.Fatal(app.Graceful(":8080"))
}
