package owl

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// App is the main DX application.
type App struct {
	mux          *Mux
	errorHandler ErrorHandler
	middlewares  []Middleware
	name         string // Server name (default: "Owl")
	version      string // Server version (default: Version constant)
	bodyLimit    int64  // Max request body size in bytes (default: 10MB)
	strictJSON   bool   // Reject JSON with unknown fields (default: false)
}

// AppConfig holds configuration for creating a new App.
type AppConfig struct {
	Name       string // Server name (default: "Owl")
	Version    string // Server version (default: owl.Version)
	BodyLimit  int64  // Max request body size in bytes (default: 10MB, 0 = unlimited)
	StrictJSON bool   // Reject JSON with unknown fields (default: false)
}

// New creates a new App with optional configuration.
func New(config ...AppConfig) *App {
	app := &App{
		mux:          NewMux(),
		errorHandler: defaultErrorHandler,
		middlewares:  []Middleware{},
		name:         "Owl",
		version:      Version,
		bodyLimit:    10 * MB, // 10MB default
		strictJSON:   false,   // Allow unknown fields by default
	} // Apply config if provided
	if len(config) > 0 {
		cfg := config[0]
		if cfg.Name != "" {
			app.name = cfg.Name
		}
		if cfg.Version != "" {
			app.version = cfg.Version
		}
		if cfg.BodyLimit > 0 {
			app.bodyLimit = cfg.BodyLimit
		} else if cfg.BodyLimit == 0 {
			// 0 means unlimited (remove limit)
			app.bodyLimit = 0
		}
		app.strictJSON = cfg.StrictJSON
	}

	return app
}

// Use adds middlewares - accepts both net/http and Owl-style.
// It automatically detects and applies the correct type:
//   - func(http.Handler) http.Handler (chi/standard middleware)
//   - func(Handler) Handler (Owl-style middleware)
func (a *App) Use(middlewares ...interface{}) *App {
	for _, mw := range middlewares {
		switch m := mw.(type) {
		case func(http.Handler) http.Handler:
			// Standard net/http middleware (chi-style)
			a.mux.Use(m)
		case Middleware:
			// Owl-style middleware
			a.middlewares = append(a.middlewares, m)
		default:
			panic("middleware must be either func(http.Handler) http.Handler or func(Handler) Handler")
		}
	}
	return a
}

// SetErrorHandler sets custom error handler.
func (a *App) SetErrorHandler(h ErrorHandler) *App {
	a.errorHandler = h
	return a
}

// Group creates a route group with prefix and middlewares.
func (a *App) Group(prefix string, middlewares ...Middleware) *Group {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(a.middlewares))
	copy(mws, a.middlewares)
	mws = append(mws, middlewares...)

	return &Group{
		app:         a,
		prefix:      prefix,
		middlewares: mws,
	}
}

// Mux returns the underlying chi Mux for advanced usage or chi-style routing.
func (a *App) Mux() *Mux {
	return a.mux
}

// ServeHTTP implements http.Handler.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// Start starts the HTTP server (blocking).
func (a *App) Start(addr string) error {
	log.Printf("\033[92m%s\033[0m v%s server starting on \033[102;30m%s\033[0m", a.name, a.version, addr)
	return http.ListenAndServe(addr, a)
}

// Graceful starts the HTTP server with graceful shutdown support.
func (a *App) Graceful(addr string, timeout ...time.Duration) error {
	// Default timeout is 10 seconds
	shutdownTimeout := 10 * time.Second
	if len(timeout) > 0 {
		shutdownTimeout = timeout[0]
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: a,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("\033[92m%s\033[0m v%s server starting on \033[102;30m%s\033[0m", a.name, a.version, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Printf("\033[92m%s\033[0m \033[33mShutting down server gracefully...\033[0m", a.name)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("\033[41m Error \033[0m Server forced to shutdown: %v", err)
		return err
	}

	log.Printf("\033[92m%s\033[0m Server stopped", a.name)
	return nil
}

// wrapHandler converts DX Handler to http.HandlerFunc.
func (a *App) wrapHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply body limit if configured
		if a.bodyLimit > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, a.bodyLimit)
		}

		c := newCtx(w, r, a.strictJSON)
		if err := h(c); err != nil {
			a.errorHandler(c, err)
		}
	}
}

// chainMiddlewares chains middlewares (pre-compiled at registration).
func chainMiddlewares(h Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
