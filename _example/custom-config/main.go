package main

import (
	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Example showing custom app configuration
func main() {
	// Create app with custom name and version
	app := owl.New(owl.AppConfig{
		Name:    "ProductAPI",
		Version: "2.1.0",
	})

	// Add middleware
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	// Routes
	api := app.Group("/api/v1")

	api.GET("/products", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"products": []map[string]string{
				{"id": "1", "name": "Laptop"},
				{"id": "2", "name": "Phone"},
			},
		})
	})

	api.GET("/products/{id}", func(c *owl.Ctx) error {
		id := c.Param("id")
		return c.JSON(map[string]string{
			"id":   id,
			"name": "Product " + id,
		})
	})

	api.GET("/health", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status":  "ok",
			"service": "ProductAPI",
			"version": "2.1.0",
		})
	})

	// Start with graceful shutdown
	app.Graceful(":8082")
}
