package main

import (
	"log"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Example showing simple Start() without graceful shutdown
func main() {
	app := owl.New()

	// Add standard middleware using UseHTTP
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.Group("/api").GET("/hello", func(c *owl.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Hello from simple server!",
		})
	})

	// Simple blocking start (no graceful shutdown)
	log.Fatal(app.Start(":8081"))
}
