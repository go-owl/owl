package main

import (
	"log"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

func main() {
	app := owl.New(owl.AppConfig{
		Name:    "CORS-API",
		Version: "1.0.0",
	})

	// Example 1: Allow all origins (default)
	app.Use(middleware.CORS())

	// Example 2: Custom CORS config
	// app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowHeaders:     []string{"Authorization", "Content-Type"},
	// 	ExposeHeaders:    []string{"X-Total-Count"},
	// 	AllowCredentials: true,
	// 	MaxAge:           3600,
	// }))

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	api := app.Group("/api")

	api.GET("/hello", func(c *owl.Ctx) error {
		return c.JSON(map[string]string{
			"message": "CORS is enabled! 🦉",
		})
	})

	api.POST("/data", func(c *owl.Ctx) error {
		var body map[string]interface{}
		if err := c.Bind().JSON(&body); err != nil {
			return err
		}

		return c.Status(201).JSON(map[string]interface{}{
			"success": true,
			"data":    body,
		})
	})

	log.Fatal(app.Graceful(":8080"))
}
