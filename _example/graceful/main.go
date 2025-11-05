package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Example showing how to implement graceful shutdown manually.
// This gives you full control over the shutdown process.
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

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8082",
		Handler: app,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on :8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server gracefully...")

	// Create shutdown context with timeout (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

	// Try: curl http://localhost:8082/api/long-task
	// Then Ctrl+C to see graceful shutdown waiting for request to complete
}
