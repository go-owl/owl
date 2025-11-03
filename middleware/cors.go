package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig defines CORS configuration.
type CORSConfig struct {
	// AllowOrigins defines allowed origins. Use ["*"] to allow all.
	AllowOrigins []string

	// AllowMethods defines allowed HTTP methods.
	AllowMethods []string

	// AllowHeaders defines allowed request headers.
	AllowHeaders []string

	// ExposeHeaders defines which headers are safe to expose.
	ExposeHeaders []string

	// AllowCredentials indicates whether credentials are allowed.
	AllowCredentials bool

	// MaxAge indicates how long preflight results can be cached (in seconds).
	MaxAge int
}

// DefaultCORSConfig returns a default CORS configuration.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           3600,
	}
}

// CORS returns a CORS middleware with default config.
func CORS() func(http.Handler) http.Handler {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig returns a CORS middleware with custom config.
func CORSWithConfig(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowOrigin := ""
			for _, o := range config.AllowOrigins {
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
			}

			// If origin not allowed and not wildcard, skip CORS headers
			if allowOrigin == "" && len(config.AllowOrigins) > 0 && config.AllowOrigins[0] != "*" {
				next.ServeHTTP(w, r)
				return
			}

			// Set origin
			if allowOrigin == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}

			// Set credentials
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight request
			if r.Method == http.MethodOptions {
				// Set allowed methods
				if len(config.AllowMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
				}

				// Set allowed headers
				if len(config.AllowHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
				} else {
					// Echo back requested headers
					h := r.Header.Get("Access-Control-Request-Headers")
					if h != "" {
						w.Header().Set("Access-Control-Allow-Headers", h)
					}
				}

				// Set max age
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Set exposed headers for actual request
			if len(config.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
			}

			next.ServeHTTP(w, r)
		})
	}
}
