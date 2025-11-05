package main

import (
	"context"
	"log"
	"time"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
	"go.uber.org/fx"
)

// User represents a user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest represents request body for creating users
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewApp creates a new Owl app with configuration
func NewApp() *owl.App {
	app := owl.New(owl.AppConfig{
		Name:      "Professional Owl API",
		Version:   "2.0.0",
		BodyLimit: 10 * owl.MB,
	})

	// Add middleware stack
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)

	return app
}

// RegisterRoutes sets up all application routes
func RegisterRoutes(app *owl.App) {
	// Health check endpoint
	app.GET("/health", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "2.0.0",
		})
	})

	// API v1 group
	v1 := app.Group("/api/v1")

	// Users endpoints
	v1.GET("/users", func(c *owl.Ctx) error {
		users := []User{
			{ID: 1, Name: "Alice Johnson", Email: "alice@example.com"},
			{ID: 2, Name: "Bob Smith", Email: "bob@example.com"},
			{ID: 3, Name: "Charlie Brown", Email: "charlie@example.com"},
		}
		return c.JSON(map[string]interface{}{
			"success": true,
			"data":    users,
			"count":   len(users),
		})
	})

	v1.POST("/users", func(c *owl.Ctx) error {
		var req CreateUserRequest

		// Use Auto binder for content-type detection
		if err := c.Bind().Auto(&req); err != nil {
			return err
		}

		// Validate required fields
		if req.Name == "" || req.Email == "" {
			return owl.NewHTTPError(400, "name and email are required")
		}

		// Create new user (simulate database)
		user := User{
			ID:    999,
			Name:  req.Name,
			Email: req.Email,
		}

		return c.Status(201).JSON(map[string]interface{}{
			"success": true,
			"message": "User created successfully",
			"data":    user,
		})
	})

	v1.GET("/users/{id}", func(c *owl.Ctx) error {
		id := c.Param("id")

		// Simulate user lookup
		user := User{
			ID:    1,
			Name:  "User " + id,
			Email: "user" + id + "@example.com",
		}

		return c.JSON(map[string]interface{}{
			"success": true,
			"data":    user,
		})
	})
}

// RegisterHTTPLifecycle manages server lifecycle with UberFx
func RegisterHTTPLifecycle(lc fx.Lifecycle, app *owl.App) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("🚀 Starting HTTP server on :8080")
			go func() {
				server := app.Listen(":8080")
				if err := server.ListenAndServe(); err != nil {
					log.Printf("❌ Server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("🛑 Shutting down HTTP server gracefully...")
			return app.Shutdown()
		},
	})
}

// Module defines the fx dependency injection module
var Module = fx.Options(
	fx.Provide(NewApp),
	fx.Invoke(RegisterRoutes),
	fx.Invoke(RegisterHTTPLifecycle),
)

func main() {
	// Create and start fx application
	app := fx.New(
		Module,
		fx.NopLogger, // Suppress fx internal logs
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	log.Println("🦉 Professional Owl API with UberFx is running!")
	log.Println("📋 Available endpoints:")
	log.Println("   GET  /health")
	log.Println("   GET  /api/v1/users")
	log.Println("   POST /api/v1/users")
	log.Println("   GET  /api/v1/users/{id}")
	log.Println("")
	log.Println("💡 Try:")
	log.Println("   curl http://localhost:8080/health")
	log.Println("   curl http://localhost:8080/api/v1/users")
	log.Println("   curl -X POST http://localhost:8080/api/v1/users -H 'Content-Type: application/json' -d '{\"name\":\"John\",\"email\":\"john@example.com\"}'")

	// Wait for shutdown signal
	<-app.Done()
	log.Println("✅ Application stopped gracefully")
}
