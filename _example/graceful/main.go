package main

import (
	"log"
	"time"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Example showing Graceful() with custom timeout
func main() {
	app := owl.New(owl.AppConfig{
		Name:    "GracefulAPI",
		Version: "1.0.0",
	})

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	app.Group("/api").GET("/long-task", func(c *owl.Ctx) error {
		// Simulate long-running task
		time.Sleep(2 * time.Second)
		return c.JSON(map[string]string{
			"message": "Task completed!",
		})
	})

	// Graceful shutdown with custom timeout (30 seconds)
	// Try: curl http://localhost:8082/api/long-task
	// Then Ctrl+C to see graceful shutdown waiting for request to complete
	log.Fatal(app.Graceful(":8082", 30*time.Second))
}
