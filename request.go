package owl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// HTTPError represents an HTTP error with code and message.
type HTTPError struct {
	Code    int
	Message string
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}

// Common HTTP errors.
var (
	ErrBadRequest   = &HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrUnauthorized = &HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden    = &HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound     = &HTTPError{Code: http.StatusNotFound, Message: "Not Found"}
)

// BindJSON decodes JSON from request body into dst.
func BindJSON(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return NewHTTPError(http.StatusBadRequest, "request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid JSON: "+err.Error())
	}
	return nil
}

// Query returns the query parameter value by key.
func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// Header returns the request header value by key.
func Header(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// ClientIP returns the client IP address.
// If trustProxy is true, checks X-Real-IP and X-Forwarded-For headers.
func ClientIP(r *http.Request, trustProxy bool) string {
	if trustProxy {
		// Check X-Real-IP
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			return ip
		}
		// Check X-Forwarded-For
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// Take the first IP
			if idx := strings.Index(xff, ","); idx > 0 {
				return strings.TrimSpace(xff[:idx])
			}
			return strings.TrimSpace(xff)
		}
	}
	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > 0 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}
