package owl

import (
	"net/http"
)

// Ctx represents the request context.
type Ctx struct {
	Request  *http.Request
	Response http.ResponseWriter
	status   int
}

// newCtx creates a new Ctx.
func newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		Request:  r,
		Response: w,
		status:   http.StatusOK,
	}
}

// Param retrieves URL path parameter.
func (c *Ctx) Param(key string) string {
	return URLParam(c.Request, key)
}

// Query retrieves URL query parameter.
func (c *Ctx) Query(key string) string {
	return Query(c.Request, key)
}

// Header retrieves request header.
func (c *Ctx) Header(key string) string {
	return Header(c.Request, key)
}

// SetHeader sets response header.
func (c *Ctx) SetHeader(key, value string) *Ctx {
	c.Response.Header().Set(key, value)
	return c
}

// Status sets response status code.
func (c *Ctx) Status(code int) *Ctx {
	c.status = code
	return c
}

// Bind returns a Binder for flexible content type binding.
// Example: c.Bind().JSON(&data), c.Bind().XML(&data)
func (c *Ctx) Bind() *Binder {
	return &Binder{
		request: c.Request,
	}
}

// BindJSON binds request JSON body to dst.
// Deprecated: Use c.Bind().JSON(dst) for more flexibility.
// This method is kept for backward compatibility.
func (c *Ctx) BindJSON(dst interface{}) error {
	return BindJSON(c.Request, dst)
}

// JSON sends JSON response.
func (c *Ctx) JSON(data interface{}) error {
	return JSON(c.Response, c.status, data)
}

// Text sends plain text response.
func (c *Ctx) Text(text string) error {
	return Text(c.Response, c.status, text)
}

// ClientIP returns client IP address.
func (c *Ctx) ClientIP(trustProxy bool) string {
	return ClientIP(c.Request, trustProxy)
}

// Handler is the DX layer handler that returns an error.
type Handler func(*Ctx) error

// Middleware wraps a Handler.
type Middleware func(Handler) Handler

// ErrorHandler handles errors from handlers.
type ErrorHandler func(*Ctx, error)

// defaultErrorHandler sends JSON error response.
func defaultErrorHandler(c *Ctx, err error) {
	if err == nil {
		return
	}

	// Check if it's an HTTPError
	if httpErr, ok := err.(*HTTPError); ok {
		_ = JSON(c.Response, httpErr.Code, map[string]interface{}{
			"success": false,
			"code":    httpErr.Code,
			"message": httpErr.Message,
		})
		return
	}

	// Unknown error -> 500
	_ = JSON(c.Response, http.StatusInternalServerError, map[string]interface{}{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": err.Error(),
	})
}
