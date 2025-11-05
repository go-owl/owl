package main

import (
	"log"
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

func main() {
	app := owl.New(owl.AppConfig{
		Name:    "PermissionAPI",
		Version: "1.0.0",
	})

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	api := app.Group("/api/v1")

	// Method-level middleware example
	api.Route("/users").
		GET(listUsers, RequirePermission("user:read")).
		POST(createUser, RequirePermission("user:write"))

	api.Route("/users/{id}").
		GET(getUser, RequirePermission("user:read")).
		PUT(updateUser, RequirePermission("user:write")).
		DELETE(deleteUser, RequirePermission("user:delete"))

	// Admin routes with multiple middleware
	api.Route("/admin/stats").
		GET(getStats, RequireAuth(), RequireAdmin(), RateLimit())

	log.Fatal(app.Start(":8080"))
}

// Permission middleware
func RequirePermission(permission string) owl.Middleware {
	return func(next owl.Handler) owl.Handler {
		return func(c *owl.Ctx) error {
			// Check permission from context/token
			userPerms := c.Header("X-Permissions")

			// Simple check (in real app, parse JWT and check permissions)
			if userPerms != permission && userPerms != "admin:*" {
				return owl.NewHTTPError(http.StatusForbidden,
					"Missing permission: "+permission)
			}

			log.Printf("✅ Permission check passed: %s", permission)
			return next(c)
		}
	}
}

// Auth middleware
func RequireAuth() owl.Middleware {
	return func(next owl.Handler) owl.Handler {
		return func(c *owl.Ctx) error {
			token := c.Header("Authorization")
			if token == "" {
				return owl.NewHTTPError(http.StatusUnauthorized, "Missing auth token")
			}
			log.Println("✅ Auth check passed")
			return next(c)
		}
	}
}

// Admin middleware
func RequireAdmin() owl.Middleware {
	return func(next owl.Handler) owl.Handler {
		return func(c *owl.Ctx) error {
			role := c.Header("X-Role")
			if role != "admin" {
				return owl.NewHTTPError(http.StatusForbidden, "Admin only")
			}
			log.Println("✅ Admin check passed")
			return next(c)
		}
	}
}

// Rate limit middleware
func RateLimit() owl.Middleware {
	return func(next owl.Handler) owl.Handler {
		return func(c *owl.Ctx) error {
			// Simplified rate limiting
			log.Println("✅ Rate limit check passed")
			return next(c)
		}
	}
}

// Handlers
func listUsers(c *owl.Ctx) error {
	return c.JSON(map[string]interface{}{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
}

func createUser(c *owl.Ctx) error {
	return c.Status(http.StatusCreated).JSON(map[string]string{
		"message": "User created",
	})
}

func getUser(c *owl.Ctx) error {
	id := c.Param("id")
	return c.JSON(map[string]string{
		"id":   id,
		"name": "User " + id,
	})
}

func updateUser(c *owl.Ctx) error {
	id := c.Param("id")
	return c.JSON(map[string]string{
		"message": "User " + id + " updated",
	})
}

func deleteUser(c *owl.Ctx) error {
	id := c.Param("id")
	return c.JSON(map[string]string{
		"message": "User " + id + " deleted",
	})
}

func getStats(c *owl.Ctx) error {
	return c.JSON(map[string]interface{}{
		"total_users":  100,
		"active_users": 85,
		"admin_only":   true,
	})
}
