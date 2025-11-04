package owl

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBodyLimit(t *testing.T) {
	tests := []struct {
		name        string
		bodyLimit   int64
		bodySize    int
		expectError bool
	}{
		{
			name:        "Within limit",
			bodyLimit:   1024,
			bodySize:    512,
			expectError: false,
		},
		{
			name:        "Exceeds limit",
			bodyLimit:   1024,
			bodySize:    2048,
			expectError: true,
		},
		{
			name:        "Unlimited (0)",
			bodyLimit:   0,
			bodySize:    10000,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := New(AppConfig{
				BodyLimit: tt.bodyLimit,
			})

			app.Group("").POST("/upload", func(c *Ctx) error {
				var data map[string]interface{}
				if err := c.BindJSON(&data); err != nil {
					return err
				}
				return c.JSON(map[string]string{"status": "ok"})
			})

			// Create request with specified body size
			body := strings.Repeat("x", tt.bodySize)
			req := httptest.NewRequest("POST", "/upload", bytes.NewBufferString(`{"data":"`+body+`"}`))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			app.ServeHTTP(w, req)

			if tt.expectError {
				if w.Code == http.StatusOK {
					t.Errorf("Expected error response, got status 200")
				}
			} else {
				if w.Code != http.StatusOK {
					t.Errorf("Expected status 200, got %d", w.Code)
				}
			}
		})
	}
}

func TestDefaultBodyLimit(t *testing.T) {
	app := New() // Should have 10MB default

	expectedLimit := 10 * MB // 10MB
	if app.bodyLimit != expectedLimit {
		t.Errorf("Expected default body limit 10MB (%d), got %d", expectedLimit, app.bodyLimit)
	}
}
