package main

import (
	"log"
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

// User represents a user model.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Create new app with custom config
	app := owl.New(owl.AppConfig{
		Name:    "HelloWorldAPI",
		Version: "1.0.0",
	})

	// Add middleware
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.RequestID)

	// API group with Auth middleware
	api := app.Group("/api")

	// v1 group with Permission middleware
	v1 := api.Group("/v1")

	// Route chaining style
	v1.Route("/users").
		GET(listUsers).
		POST(createUser)

	v1.Route("/users/{id}").
		GET(getUser).
		PATCH(updateUser).
		DELETE(deleteUser)

	// Additional routes
	v1.GET("/health", healthCheck)
	v1.GET("/echo", echoQuery)

	// Start server
	log.Fatal(app.Start(":8080"))
}

// Handler functions

func healthCheck(c *owl.Ctx) error {
	return c.JSON(map[string]interface{}{
		"status":  "ok",
		"message": "Owl is flying high!",
	})
}

func echoQuery(c *owl.Ctx) error {
	name := c.Query("name")
	if name == "" {
		name = "anonymous"
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "Hello, " + name + "!",
		"query":   c.Request.URL.RawQuery,
	})
}

func listUsers(c *owl.Ctx) error {
	// Simulate database query
	users := []User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
		{ID: "3", Name: "Charlie", Email: "charlie@example.com"},
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"success": true,
		"data":    users,
		"count":   len(users),
	})
}

func createUser(c *owl.Ctx) error {
	var req CreateUserRequest

	// Bind JSON from request body (new flexible API)
	if err := c.Bind().JSON(&req); err != nil {
		return err // Will be handled by error handler
	}

	// Validate
	if req.Name == "" {
		return owl.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	if req.Email == "" {
		return owl.NewHTTPError(http.StatusBadRequest, "email is required")
	}

	// Simulate creating user
	user := User{
		ID:    "999",
		Name:  req.Name,
		Email: req.Email,
	}

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
}

func getUser(c *owl.Ctx) error {
	userID := c.Param("id")

	// Simulate database query
	if userID == "" {
		return owl.NewHTTPError(http.StatusBadRequest, "user ID is required")
	}

	user := User{
		ID:    userID,
		Name:  "User " + userID,
		Email: "user" + userID + "@example.com",
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func updateUser(c *owl.Ctx) error {
	userID := c.Param("id")

	var req CreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	// Simulate updating user
	user := User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "User updated successfully",
	})
}

func deleteUser(c *owl.Ctx) error {
	userID := c.Param("id")

	// Simulate deleting user
	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"success": true,
		"message": "User " + userID + " deleted successfully",
	})
}
