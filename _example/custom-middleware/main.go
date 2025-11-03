package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Custom context key for storing user info
type contextKey string

const userContextKey = contextKey("user")

// Example showing how to create custom middleware for Owl
func main() {
	app := owl.New(owl.AppConfig{
		Name:    "CustomMiddlewareAPI",
		Version: "1.0.0",
	})

	// Add built-in middleware
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	// Add custom chi-style middleware (works with app.Use)
	app.Use(timingMiddleware)
	app.Use(authMiddleware)

	// Public routes (no auth)
	root := app.Group("")
	root.GET("/", homeHandler)
	root.GET("/health", healthHandler)

	// Protected routes (requires auth)
	api := app.Group("/api")
	api.GET("/profile", profileHandler)
	api.GET("/data", dataHandler)

	log.Fatal(app.Graceful(":8083"))
}

// timingMiddleware measures request processing time (chi-style)
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Process request
		next.ServeHTTP(w, r)

		// Log processing time
		duration := time.Since(start)
		log.Printf("⏱️  Request took %v", duration)
	})
}

// authMiddleware checks for API key (chi-style)
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public routes
		if r.URL.Path == "/" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Check API key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, `{"error":"API key required"}`, http.StatusUnauthorized)
			return
		}

		if apiKey != "secret123" {
			http.Error(w, `{"error":"Invalid API key"}`, http.StatusForbidden)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), userContextKey, "authenticated-user")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handlers

func homeHandler(c *owl.Ctx) error {
	return c.JSON(map[string]interface{}{
		"message": "Welcome! This is a public route.",
		"tip":     "Try /api/profile with X-API-Key: secret123",
	})
}

func healthHandler(c *owl.Ctx) error {
	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func profileHandler(c *owl.Ctx) error {
	// Get user from context
	user := c.Request.Context().Value(userContextKey)

	return c.JSON(map[string]interface{}{
		"message": "Protected profile data",
		"user":    user,
		"data": map[string]string{
			"name":  "John Doe",
			"email": "john@example.com",
		},
	})
}

func dataHandler(c *owl.Ctx) error {
	return c.JSON(map[string]interface{}{
		"message": "Protected data endpoint",
		"data":    []int{1, 2, 3, 4, 5},
	})
}
