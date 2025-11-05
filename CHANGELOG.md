# Changelog

All notable changes to Owl will be documented in this file.

## [Unreleased]

### Added

- 🔒 **StrictJSON Mode**: Production-ready JSON validation
  - `AppConfig.StrictJSON` - Reject unknown fields via `DisallowUnknownFields()`
  - Detect and reject trailing data after JSON objects
  - Prevent typos, injection attacks, and API contract violations
- 🤖 **Auto Binder**: `c.Bind().Auto(&data)` - Automatic content-type detection
  - Supports JSON, XML, Form, and Multipart automatically
  - Eliminates manual Content-Type checking
- 📊 **Advanced Slice Binding**: Support all primitive types in slices
  - `[]string`, `[]int`, `[]int64`, `[]float64`, `[]bool`
  - Example: `?tags=a&tags=b&scores=1&scores=2`
- 🎯 **Pointer & Array Field Support**:
  - Bind to pointer fields: `*string`, `*int`
  - Bind to array fields: `[3]int`, `[5]string`
- 🚀 **HTTP Method Shortcuts**: Direct routing on App instance
  - `app.GET()`, `app.POST()`, `app.PUT()`, `app.PATCH()`, `app.DELETE()`
  - Convenience methods for simple routes without groups
  - Maintains full middleware and handler chaining support
- 🔌 **UberFx Compatibility**: Simple lifecycle management
  - `app.Listen(addr)` returns `*http.Server` for external management
  - `app.Shutdown()` for graceful shutdown (similar to Fiber)
  - Perfect for dependency injection frameworks
- 🛡️ **Enhanced Security**:
  - Named constants: `maxFieldLength`, `maxTextBodySize`, `maxFileSize`
  - File size validation (50MB per file)
  - Field length limits (10KB per field)
- ♻️ **Code Quality Improvements**:
  - DRY refactor: `readBodyLimited()` helper for Text/Bytes methods
  - Eliminated code duplication
  - Professional constant naming
- 📚 **Reorganized Examples**: 8 comprehensive examples
  - `rest-api/` - Complete CRUD API
  - `request-binding/` - All binding methods
  - `strict-json/` - StrictJSON validation demo
  - `graceful-shutdown/` - Production shutdown
  - `cors/` - CORS configuration
  - `request-limits/` - Body size limits
  - `middleware-chain/` - Custom middleware
  - `method-level-middleware/` - Per-method middleware
  - `hybrid-routing/` - Mixed Owl/chi styles
- 📖 **Documentation Updates**:
  - Updated `_example/README.md` with all examples
  - Added learning path for beginners and advanced users
  - Comprehensive curl examples

### Changed

- 🔧 Refactored `Text()` and `Bytes()` methods to use shared helper
- 🔧 Replaced magic numbers with named constants throughout codebase
- 🔧 Improved error messages for better debugging

### Removed

- 🗑️ **Removed Graceful() Method**: Keep framework focused
  - Removed `app.Graceful()` method for simplicity
  - Added simple `app.Listen()` and `app.Shutdown()` for uberfx compatibility
  - Similar API to Fiber framework for familiar experience
- 🗑️ **Removed StrictJSON Feature**: Simplified JSON binding
  - Removed `AppConfig.StrictJSON` and related validation logic
  - Developers can use `json.Decoder.DisallowUnknownFields()` directly if needed
  - Reduces framework size and keeps it focused on core features
  - Other major frameworks (Gin, Echo, Fiber) don't include this feature

### Fixed

- 🔒 **Security Fix**: Eliminated redundant body limit checks
  - Removed duplicate `maxTextBodySize` constant in binder
  - Now relies solely on App-level `MaxBytesReader` for consistent protection
  - Prevents conflicting body limit behaviors

### Improved

- ✅ All 118 tests passing
- ✅ No performance regression
- ✅ Better code maintainability
- ✅ Enhanced type safety

## [1.0.0] - 2025-11-04

### Added - Owl Fork

**Express-style Layer on chi v5:**

- ✨ **App API**: `owl.New()` for creating Express-style applications
- ✨ **Context API**: `owl.Ctx` wrapping request/response with Express-like methods:
  - `c.Param(key)` - Get URL parameters
  - `c.Query(key)` - Get query parameters
  - `c.Header(key)` - Get request headers
  - `c.SetHeader(key, value)` - Set response headers
  - `c.Status(code)` - Set response status
  - `c.JSON(data)` - Send JSON response
  - `c.Text(text)` - Send text response
  - `c.BindJSON(&dst)` - Parse JSON request body
- ✨ **Route Chaining**: `app.Route("/users").GET(h1).POST(h2).PUT(h3)`
- ✨ **Route Groups**: `app.Group("/api").GET("/users", handler)`
- ✨ **Error Handling**:
  - Handlers return errors: `func(c *owl.Ctx) error`
  - `owl.NewHTTPError(code, message)` for HTTP errors
  - `app.SetErrorHandler()` for custom error handling
- ✨ **Graceful Shutdown**: `app.Graceful(":3000", timeout)` with signal handling
- ✨ **Middleware Support**: Compatible with both chi and Owl-style middleware
- 📚 **Examples**: Added Express-style examples in `_example/`

**Repository Changes:**

- 🔄 Renamed from `chi` to `owl`
- 🔄 Changed module path from `github.com/go-chi/chi/v5` to `owl`
- 🔄 Updated all imports and package names
- 📝 Created new documentation structure

### Maintained from chi v5

- ✅ **100% chi v5 API compatibility** via `app.Mux()`
- ✅ **Radix tree router** - Same performance as chi
- ✅ **All middleware included** - Complete `owl/middleware` package with 20+ production-ready middleware:
  - Core: Logger, Recoverer, RequestID, RealIP, Compress, Timeout
  - Security: BasicAuth, AllowContentType, NoCache
  - Performance: Throttle, ContentEncoding
  - All chi middleware work out-of-the-box
- ✅ **URL patterns** - Named params `{id}`, wildcards `*`, regex `{id:\\d+}`
- ✅ **Context values** - Request-scoped values via `context.Context`
- ✅ **Standard handlers** - Full `net/http` compatibility
- ✅ **External packages** - Compatible with chi ecosystem (cors, jwtauth, httprate, etc.)

---

## Original chi v5 Changelog

For the complete chi v5 history, see: https://github.com/go-chi/chi/blob/master/CHANGELOG.md

### Notable chi v5 Features (Inherited)

**v5.0.12 (2024-02-16)**

- Latest stable chi v5 release
- See full history: https://github.com/go-chi/chi/blob/master/CHANGELOG.md

**v5.0.0 (2021-02-27)**

- Introduced Semantic Import Versioning (SIV)
- Full Go modules support
- Context integration improvements
- See full history: https://github.com/go-chi/chi/compare/v1.5.4...v5.0.0

---

## Credits

Owl is built on [go-chi/chi](https://github.com/go-chi/chi) v5.0.12

**Original chi authors:**

- Peter Kieltyka ([@pkieltyka](https://github.com/pkieltyka))
- Vojtech Vitek ([@VojtechVitek](https://github.com/VojtechVitek))

**Influenced by:**

- Express.js - API design inspiration
- Fiber - Context and chaining patterns
- goji - Middleware patterns (Carl Jackson)
- go-radix - Radix tree implementation (Armon Dadgar)

---

## License

MIT License

- Owl: Copyright (c) 2025 Owl Contributors
- chi: Copyright (c) 2015-present Peter Kieltyka

---
