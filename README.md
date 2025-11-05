# Owl 🦉

**Express-style HTTP router built on [go-chi/chi](https://github.com/go-chi/chi)**

Lightweight, fast, and idiomatic Go HTTP router with Express.js-inspired API.

---

## Features

- 🚀 **Express-like API** - `app.Get()`, `c.JSON()`, `c.Param()` style
- ⚡ **chi Performance** - Battle-tested radix tree routing
- 🔧 **100% Compatible** - Works with any `net/http` handler
- 🎨 **Method Chaining** - Clean, fluent route definitions
- 🛡️ **Error Handling** - Built-in error propagation
- 🔒 **Advanced Binding** - JSON, XML, Form, Query, Multipart with StrictJSON mode
- 🤖 **Auto Binder** - Automatic content-type detection
- 🌳 **Zero Dependencies** - Pure standard library

## Quick Start

```go
package main

import "github.com/go-owl/owl"

func main() {
    app := owl.New()

    app.GET("/hello", func(c *owl.Ctx) error {
        return c.JSON(map[string]string{"message": "Hello, Owl! 🦉"})
    })

    app.Graceful(":3000")
}
```

## Installation

```bash
go get github.com/go-owl/owl
```

Or add to your `go.mod`:

```go
require github.com/go-owl/owl v1.0.0
```

## Examples

### CRUD API

```go
import "log"

func main() {
    app := owl.New()
    api := app.Group("/api/v1")

    // Method chaining
    api.Route("/users").
        GET(listUsers).
        POST(createUser)

    api.Route("/users/{id}").
        GET(getUser).
        PUT(updateUser).
        DELETE(deleteUser)

    log.Fatal(app.Start(":3000"))
}

func listUsers(c *owl.Ctx) error {
    users := []User{{ID: "1", Name: "Alice"}}
    return c.JSON(users)
}

func getUser(c *owl.Ctx) error {
    id := c.Param("id")
    return c.JSON(User{ID: id, Name: "User " + id})
}
```

### Middleware

```go
import "github.com/go-owl/owl/middleware"

app := owl.New()

// Standard middleware (inherited from chi)
app.Use(middleware.Logger)
app.Use(middleware.Recoverer)
app.Use(middleware.RequestID)
app.Use(middleware.RealIP)
app.Use(middleware.Compress(5))

// Custom Owl-style middleware
func Auth(next owl.Handler) owl.Handler {
    return func(c *owl.Ctx) error {
        if c.Header("Authorization") == "" {
            return owl.NewHTTPError(401, "Unauthorized")
        }
        return next(c)
    }
}

app.Group("/api", Auth).GET("/protected", handler)
```

> **Note:** Owl includes **all chi middleware** - battle-tested in production by Cloudflare, Heroku, and thousands of projects. Use `net/http` compatible middleware from chi ecosystem too!

### Request Binding

```go
import "mime/multipart"

// Enable StrictJSON for production (rejects unknown fields & trailing data)
app := owl.New(owl.AppConfig{
    StrictJSON: true,
    BodyLimit:  10 * owl.MB,
})

func createUser(c *owl.Ctx) error {
    var user User

    // Flexible binding methods
    if err := c.Bind().JSON(&user); err != nil {
        return err
    }

    // Or use Auto() - detects content type automatically
    if err := c.Bind().Auto(&user); err != nil {
        return err
    }

    return c.Status(201).JSON(user)
}

func search(c *owl.Ctx) error {
    var query struct {
        Q     string   `query:"q"`
        Tags  []string `query:"tags"`  // Supports arrays: ?tags=a&tags=b
        Page  int      `query:"page"`
    }

    c.Bind().Query(&query)
    return c.JSON(query)
}

func upload(c *owl.Ctx) error {
    var form struct {
        Title string                `form:"title"`
        File  *multipart.FileHeader `form:"file"`
    }

    c.Bind().MultipartForm(&form, 10*owl.MB)
    return c.JSON(map[string]interface{}{
        "filename": form.File.Filename,
        "size":     form.File.Size,
    })
}
```

**Supported bindings:**

- `c.Bind().JSON(&dst)` - JSON with optional StrictJSON mode
- `c.Bind().XML(&dst)` - XML parsing
- `c.Bind().Form(&dst)` - URL-encoded forms
- `c.Bind().Query(&dst)` - Query parameters with array support
- `c.Bind().MultipartForm(&dst, maxMemory)` - File uploads
- `c.Bind().Text(&str)` - Raw text (webhooks)
- `c.Bind().Bytes(&bytes)` - Raw bytes
- `c.Bind().Auto(&dst)` - Auto-detect content type

## API Highlights

### Context

```go
func handler(c *owl.Ctx) error {
    // Request
    id := c.Param("id")                    // Path params
    name := c.Query("name")                // Query params
    token := c.Header("Authorization")     // Headers

    // Flexible binding (new style)
    var body User
    c.Bind().JSON(&body)                   // Parse JSON
    c.Bind().XML(&body)                    // Parse XML
    c.Bind().Form(&body)                   // Parse form data
    c.Bind().Query(&body)                  // Parse query params
    c.Bind().Auto(&body)                   // Auto-detect content type

    // Response
    return c.Status(200).JSON(body)
    return c.Text("Hello")
    return c.SetHeader("X-Custom", "value").JSON(data)
}
```

### Routing

```go
app := owl.New()

// Simple routes
app.GET("/", home)
app.POST("/users", createUser)

// URL parameters
app.GET("/users/{id}", getUser)
app.GET("/posts/{slug:[a-z-]+}", getPost) // Regex

// Groups
api := app.Group("/api/v1")
api.GET("/health", healthCheck)

// Chi-style (also supported)
app.Mux().Get("/chi-style", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("traditional chi handler"))
})
```

## Architecture

Owl is a thin Express-style wrapper around chi's router:

```
┌─────────────────────────────┐
│  Owl Express Layer          │ ← app.GET, c.JSON, c.Param
├─────────────────────────────┤
│  chi Router (v5)            │ ← Radix tree, routing
├─────────────────────────────┤
│  Go net/http                │ ← Standard library
└─────────────────────────────┘
```

- **Owl layer** provides Express-style API
- **chi core** handles all routing (zero overhead)
- **net/http** compatibility maintained

## Credits

Built on **[go-chi/chi](https://github.com/go-chi/chi)** v5

**Original chi authors:**

- Peter Kieltyka ([@pkieltyka](https://github.com/pkieltyka))
- Vojtech Vitek ([@VojtechVitek](https://github.com/VojtechVitek))
- All [chi contributors](https://github.com/go-chi/chi/graphs/contributors)

**Additional credits:**

- Carl Jackson for [goji](https://github.com/zenazn/goji) (middleware inspiration)
- Armon Dadgar for [go-radix](https://github.com/armon/go-radix)

**Express-style enhancements** by this project.

## License

MIT License

- **Owl**: Copyright (c) 2025 Owl Contributors
- **chi**: Copyright (c) 2015-present Peter Kieltyka

See [LICENSE](./LICENSE) for details.

---

**Documentation:** [Examples](./_example) • [Migration Guide](./OWL_MIGRATION.md) • [chi Docs](https://github.com/go-chi/chi)
