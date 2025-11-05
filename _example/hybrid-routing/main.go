package main

import (
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// Example showing both Express-style and chi-style APIs
func main() {
	// Custom configuration
	app := owl.New(owl.AppConfig{
		Name:    "MixedStyleAPI",
		Version: "1.0.0",
	})

	// Add middleware (works with both styles)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	// Express-style API (Owl)
	app.Group("/api").GET("/express", expressHandler)

	// chi-style API (Traditional Go)
	app.Mux().Get("/api/chi", chiHandler)

	// Mix middlewares - use chi's Group
	app.Mux().Group(func(r owl.Router) {
		r.Get("/api/traditional", traditionalHandler)
	})

	// You can choose your preferred style!
	app.Start(":3000")
}

// Express-style handler (returns error)
func expressHandler(c *owl.Ctx) error {
	name := c.Query("name")
	if name == "" {
		name = "World"
	}

	return c.JSON(map[string]interface{}{
		"style":   "express",
		"message": "Hello, " + name + "! 🦉",
		"method":  c.Request.Method,
	})
}

// chi-style handler (standard net/http)
func chiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"style":"chi","message":"Traditional Go style"}`))
}

// Traditional handler
func traditionalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This uses chi's traditional routing"))
}
