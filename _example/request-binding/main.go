package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/go-owl/owl"
	"github.com/go-owl/owl/middleware"
)

type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
	Age   int    `json:"age" xml:"age"`
}

func main() {
	app := owl.New(owl.AppConfig{
		Name:    "BindExampleAPI",
		Version: "1.0.0",
	})

	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	api := app.Group("/api")

	// New flexible binding style
	api.POST("/users/json", func(c *owl.Ctx) error {
		var user User

		// Use new style: c.Bind().JSON()
		if err := c.Bind().JSON(&user); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "User created with new Bind().JSON()",
			"user":    user,
		})
	})

	// XML binding example
	api.POST("/users/xml", func(c *owl.Ctx) error {
		var user User

		// XML binding is also supported!
		if err := c.Bind().XML(&user); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "User created from XML",
			"user":    user,
		})
	})

	// Old style - still works for backward compatibility
	api.POST("/users/legacy", func(c *owl.Ctx) error {
		var user User

		// Legacy style still works
		if err := c.BindJSON(&user); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success": true,
			"message": "User created with legacy BindJSON()",
			"user":    user,
		})
	})

	// Text binding for webhooks
	api.POST("/webhook/text", func(c *owl.Ctx) error {
		var payload string

		// Receive raw text - perfect for webhooks
		if err := c.Bind().Text(&payload); err != nil {
			return err
		}

		return c.JSON(map[string]interface{}{
			"success": true,
			"message": "Webhook received",
			"payload": payload,
			"length":  len(payload),
		})
	})

	// Webhook with signature verification (like Stripe, GitHub)
	api.POST("/webhook/secure", func(c *owl.Ctx) error {
		var body []byte

		// Receive raw bytes for signature verification
		if err := c.Bind().Bytes(&body); err != nil {
			return err
		}

		// Verify signature (example - use environment variable in production)
		secret := "your-webhook-secret" // TODO: Use os.Getenv("WEBHOOK_SECRET") in production
		signature := c.Header("X-Signature")

		// Calculate HMAC
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(body)
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		// Use timing-safe comparison to prevent timing attacks
		if len(signature) != len(expectedSignature) || !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			return owl.NewHTTPError(http.StatusUnauthorized, "Invalid signature")
		}

		return c.JSON(map[string]interface{}{
			"success":  true,
			"message":  "Webhook verified and processed",
			"bodySize": len(body),
		})
	})

	// Query binding example
	api.GET("/search", func(c *owl.Ctx) error {
		var query struct {
			Q     string   `query:"q"`
			Page  int      `query:"page"`
			Limit int      `query:"limit"`
			Tags  []string `query:"tags"` // Support multiple values: ?tags=a&tags=b
		}

		// Bind from query parameters
		if err := c.Bind().Query(&query); err != nil {
			return err
		}

		return c.JSON(map[string]interface{}{
			"success": true,
			"query":   query.Q,
			"page":    query.Page,
			"limit":   query.Limit,
			"tags":    query.Tags, // Array of tags
		})
	})

	// Form binding example
	api.POST("/login", func(c *owl.Ctx) error {
		var credentials struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}

		// Bind from URL-encoded form
		if err := c.Bind().Form(&credentials); err != nil {
			return err
		}

		return c.JSON(map[string]interface{}{
			"success":  true,
			"username": credentials.Username,
			"message":  "Login form received",
		})
	})

	// File upload with multipart form
	api.POST("/upload", func(c *owl.Ctx) error {
		var upload struct {
			Title       string                `form:"title"`
			Description string                `form:"description"`
			File        *multipart.FileHeader `form:"file"`
		}

		// Bind multipart form including file upload
		if err := c.Bind().MultipartForm(&upload, 10*owl.MB); err != nil {
			return err
		}

		return c.Status(http.StatusCreated).JSON(map[string]interface{}{
			"success":     true,
			"title":       upload.Title,
			"description": upload.Description,
			"filename":    upload.File.Filename,
			"size":        upload.File.Size,
		})
	})

	// Health check
	api.GET("/health", func(c *owl.Ctx) error {
		return c.JSON(map[string]interface{}{
			"status":  "ok",
			"message": "Supports JSON, XML, Query, Form, MultipartForm, Text, Bytes binding",
		})
	})

	log.Println("Try these examples:")
	log.Println()
	log.Println("1. New JSON binding:")
	log.Println("   curl -X POST http://localhost:8080/api/users/json \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"age\":25}'")
	log.Println()
	log.Println("2. XML binding:")
	log.Println("   curl -X POST http://localhost:8080/api/users/xml \\")
	log.Println("     -H 'Content-Type: application/xml' \\")
	log.Println("     -d '<User><name>Jane</name><email>jane@example.com</email><age>30</age></User>'")
	log.Println()
	log.Println("3. Legacy JSON binding (still works):")
	log.Println("   curl -X POST http://localhost:8080/api/users/legacy \\")
	log.Println("     -H 'Content-Type: application/json' \\")
	log.Println("     -d '{\"name\":\"Bob\",\"email\":\"bob@example.com\",\"age\":35}'")
	log.Println()
	log.Println("4. Text binding (webhook):")
	log.Println("   curl -X POST http://localhost:8080/api/webhook/text \\")
	log.Println("     -H 'Content-Type: text/plain' \\")
	log.Println("     -d 'event=payment.success&amount=1000'")
	log.Println()
	log.Println("5. Secure webhook with signature:")
	log.Println("   PAYLOAD='event=test'")
	log.Println("   SIG=$(echo -n \"$PAYLOAD\" | openssl dgst -sha256 -hmac 'your-webhook-secret' | cut -d' ' -f2)")
	log.Println("   curl -X POST http://localhost:8080/api/webhook/secure \\")
	log.Println("     -H 'X-Signature: '$SIG \\")
	log.Println("     -d \"$PAYLOAD\"")
	log.Println()
	log.Println("6. Query binding (supports arrays):")
	log.Println("   curl 'http://localhost:8080/api/search?q=golang&page=1&limit=10&tags=api&tags=rest'")
	log.Println()
	log.Println("7. Form binding:")
	log.Println("   curl -X POST http://localhost:8080/api/login \\")
	log.Println("     -H 'Content-Type: application/x-www-form-urlencoded' \\")
	log.Println("     -d 'username=admin&password=secret'")
	log.Println()
	log.Println("8. File upload:")
	log.Println("   curl -X POST http://localhost:8080/api/upload \\")
	log.Println("     -F 'title=My Document' \\")
	log.Println("     -F 'description=Important file' \\")
	log.Println("     -F 'file=@/path/to/file.pdf'")
	log.Println()

	log.Fatal(app.Graceful(":8080"))
}
