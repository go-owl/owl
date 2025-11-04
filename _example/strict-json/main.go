package main

import (
	"log"
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func main() {
	// 🔒 Production-ready config with StrictJSON enabled
	app := owl.New(owl.AppConfig{
		Name:       "StrictJSON API",
		Version:    "1.0.0",
		BodyLimit:  10 * owl.MB,
		StrictJSON: true, // 🔥 Enable strict mode!
	})

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	api := app.Group("/api")

	// ✅ Example 1: Valid request (will succeed)
	api.POST("/users/valid", func(c *owl.Ctx) error {
		var user User

		// This will work - all fields match the struct
		if err := c.Bind().JSON(&user); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "User created successfully",
			"user":    user,
		})
	})

	// ❌ Example 2: Invalid request with unknown fields (will fail)
	api.POST("/users/unknown-fields", func(c *owl.Ctx) error {
		var user User

		// StrictJSON=true will reject unknown fields like "extraField"
		if err := c.Bind().JSON(&user); err != nil {
			// Returns: "invalid JSON: json: unknown field \"extraField\""
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"user":    user,
		})
	})

	// ❌ Example 3: Request with trailing data (will fail)
	api.POST("/users/trailing-data", func(c *owl.Ctx) error {
		var user User

		// StrictJSON=true will reject trailing data after JSON object
		if err := c.Bind().JSON(&user); err != nil {
			// Returns: "invalid JSON: trailing data after JSON object"
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"user":    user,
		})
	})

	// 📊 Example 4: Multiple endpoints with different strictness
	api.POST("/products/strict", func(c *owl.Ctx) error {
		var product Product

		// Using app-level StrictJSON config
		if err := c.Bind().JSON(&product); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "Product created with strict validation",
			"product": product,
		})
	})

	// 🔍 Example 5: Auto binder (inherits StrictJSON setting)
	api.POST("/auto", func(c *owl.Ctx) error {
		var user User

		// Auto binder also respects StrictJSON config
		if err := c.Bind().Auto(&user); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "Auto binding with strict mode",
			"user":    user,
		})
	})

	// Health check
	api.GET("/health", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status":     "ok",
			"strictJSON": true,
			"message":    "API running in strict JSON mode",
		})
	})

	log.Println("🦉 Owl API with StrictJSON enabled")
	log.Println()
	log.Println("Try these examples:")
	log.Println()
	log.Println("✅ Valid request (will succeed):")
	log.Println("   curl -X POST http://localhost:8080/api/users/valid \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"age\":25}'")
	log.Println()
	log.Println("❌ Invalid - Unknown field (will fail):")
	log.Println("   curl -X POST http://localhost:8080/api/users/unknown-fields \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"age\":25,\"extraField\":\"not allowed\"}'")
	log.Println()
	log.Println("❌ Invalid - Trailing data (will fail):")
	log.Println("   curl -X POST http://localhost:8080/api/users/trailing-data \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"age\":25}{\"extra\":\"object\"}'")
	log.Println()
	log.Println("✅ Product creation with strict validation:")
	log.Println("   curl -X POST http://localhost:8080/api/products/strict \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"Laptop\",\"price\":999.99,\"stock\":10}'")
	log.Println()
	log.Println("✅ Auto binder with strict mode:")
	log.Println("   curl -X POST http://localhost:8080/api/auto \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"Jane\",\"email\":\"jane@example.com\",\"age\":30}'")
	log.Println()
	log.Println("📝 Note: StrictJSON mode helps prevent:")
	log.Println("   - Typos in field names")
	log.Println("   - Injection of malicious fields")
	log.Println("   - API contract violations")
	log.Println("   - Trailing data attacks")
	log.Println()

	log.Fatal(app.Graceful(":8080"))
}
