# Owl Examples

Professional examples demonstrating Owl framework features.

## 📚 Available Examples

### 🔌 REST API (`rest-api/`)

Complete CRUD API with best practices:

- RESTful endpoints (GET, POST, PATCH, DELETE)
- Request/response handling with `c.Bind().JSON()`
- Path parameters and query strings
- Route groups and chaining
- Error handling

```bash
go run _example/rest-api/main.go
curl http://localhost:8080/api/v1/users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'
```

### � Request Binding (`request-binding/`)

Comprehensive binding examples for all content types:

- **JSON**: `c.Bind().JSON(&data)`
- **XML**: `c.Bind().XML(&data)`
- **Form**: `c.Bind().Form(&data)` (URL-encoded)
- **Query**: `c.Bind().Query(&data)` (URL parameters)
- **Multipart**: `c.Bind().MultipartForm(&data, maxMemory)` (file uploads)
- **Text**: `c.Bind().Text(&str)` (webhooks)
- **Bytes**: `c.Bind().Bytes(&bytes)` (raw data)

```bash
go run _example/request-binding/main.go
# See output for curl examples
```

### 🛡️ CORS (`cors/`)

Cross-Origin Resource Sharing configuration:

- Default CORS settings
- Custom CORS with specific origins
- Preflight handling

```bash
go run _example/cors/main.go
curl -X OPTIONS http://localhost:8080/api/data \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST"
```

### � Request Limits (`request-limits/`)

Body size limiting with readable constants:

- Default limit (10MB)
- Custom limits (1MB, 50MB)
- Unlimited (for special cases)
- Using `owl.KB`, `owl.MB`, `owl.GB` constants

```bash
go run _example/request-limits/main.go
curl -X POST http://localhost:8080/api/upload \
  -H "Content-Type: application/json" \
  -d '{"data":"..."}'
```

### � Middleware Chain (`middleware-chain/`)

Custom middleware implementation:

- Permission-based routing
- Request timing
- Authentication
- Context value passing

```bash
go run _example/middleware-chain/main.go
curl http://localhost:8083/api/users
curl http://localhost:8083/api/admin/settings
```

### 🎯 Method-Level Middleware (`method-level-middleware/`)

Apply middleware to specific HTTP methods:

- Different middleware for GET vs POST
- Permission checks per method
- Route-specific handlers

```bash
go run _example/method-level-middleware/main.go
curl http://localhost:8080/api/posts
curl -X POST http://localhost:8080/api/posts
```

### � Hybrid Routing (`hybrid-routing/`)

Use Express-style and chi-style APIs together:

- Owl handlers with `c *Ctx`
- Chi handlers with `w, r`
- Mixed middleware

```bash
go run _example/hybrid-routing/main.go
curl http://localhost:3000/api/owl-style
curl http://localhost:3000/api/chi-style
```

### 🛑 Graceful Shutdown (`graceful/`)

Production-ready graceful shutdown implementation:

- Signal handling (SIGINT, SIGTERM)
- Custom timeout configuration
- In-flight request completion
- Manual implementation for full control

```bash
go run _example/graceful/main.go
# Press Ctrl+C to test graceful shutdown
```

### � UberFx Integration (`uberfx/`)

Professional dependency injection with UberFx:

- Fiber-like API patterns
- Lifecycle management (OnStart/OnStop)
- Graceful shutdown with DI
- Production-ready structure

```bash
go run _example/uberfx/main.go
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
```

## 🚀 Quick Start

```bash
# Run any example
go run _example/<folder-name>/main.go

# Or build first
cd _example/<folder-name>
go build && ./<folder-name>
```

## 📖 Learning Path

**For beginners:**

1. `rest-api/` - Start here for basic CRUD operations
2. `request-binding/` - Learn all binding methods
3. `graceful/` - Production-ready graceful shutdown

**For advanced users:**

1. `middleware-chain/` - Custom middleware
2. `method-level-middleware/` - Fine-grained control
3. `hybrid-routing/` - Mix Owl and chi styles
4. `uberfx/` - Professional dependency injection

**For specific features:**

- `cors/` - Enable cross-origin requests
- `request-limits/` - Protect against large payloads

## 🎯 Key Concepts

### New Binding API

```go
// Old (deprecated but still works)
c.BindJSON(&data)

// New (recommended)
c.Bind().JSON(&data)
c.Bind().XML(&data)
c.Bind().Form(&data)
c.Bind().Query(&data)
c.Bind().MultipartForm(&data, 10*owl.MB)
```

### Production Deployment

- Always use `app.Graceful()` for production
- Set appropriate `BodyLimit` in `AppConfig`
- Add `middleware.Logger` and `middleware.Recoverer`
- Configure CORS for cross-origin APIs

## 📚 Documentation

See the main [README.md](../README.md) for complete API documentation.
