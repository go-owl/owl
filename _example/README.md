# Owl Examples

This directory contains example applications demonstrating various features of the Owl framework.

## Examples

### 🌍 Hello World

**Path:** `_example/helloworld/`

Complete example showing all basic features:

- Route groups
- Path parameters
- Query parameters
- JSON request/response
- CRUD operations
- Method chaining

```bash
go run _example/helloworld/main.go
curl http://localhost:8080/api/v1/users
```

### 🚀 Simple Start

**Path:** `_example/simple/`

Minimal example using `app.Start()` for simple HTTP server without graceful shutdown.

```bash
go run _example/simple/main.go
curl http://localhost:8081/api/hello
```

### 🛑 Graceful Shutdown

**Path:** `_example/graceful/`

Example demonstrating `app.Graceful()` with custom timeout for graceful shutdown handling.

```bash
go run _example/graceful/main.go
# In another terminal:
curl http://localhost:8082/api/long-task
# Press Ctrl+C to trigger graceful shutdown
```

### 🔀 Mixed Style

**Path:** `_example/mixed-style/`

Example showing how to use both Express-style and chi-style APIs together in one application.

```bash
go run _example/mixed-style/main.go
curl http://localhost:3000/api/express
curl http://localhost:3000/api/chi
```

### ⚙️ Custom Config

**Path:** `_example/custom-config/`

Example demonstrating custom app configuration with `AppConfig` (name and version).

```bash
go run _example/custom-config/main.go
curl http://localhost:8084/products
curl http://localhost:8084/health
```

### 🔧 Custom Middleware

**Path:** `_example/custom-middleware/`

Example showing how to create custom Owl-style middleware for:

- Request timing
- Authentication
- Context value passing

```bash
go run _example/custom-middleware/main.go
curl http://localhost:8083/
curl -H "X-API-Key: secret123" http://localhost:8083/api/profile
```

## Running Examples

Each example is a standalone application in its own directory:

```bash
# Run any example
go run _example/<example-name>/main.go

# Or build and run
cd _example/<example-name>
go build
./<example-name>
```

## Key Differences

| Feature           | `app.Start()`        | `app.Graceful()`        |
| ----------------- | -------------------- | ----------------------- |
| Server startup    | ✅ Simple blocking   | ✅ With signal handling |
| Graceful shutdown | ❌ No                | ✅ Yes                  |
| Custom timeout    | ❌ N/A               | ✅ Optional parameter   |
| Use case          | Simple apps, testing | Production apps         |

## Learn More

See the main [README.md](../README.md) for full documentation.
