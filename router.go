package owl

// Group represents a route group.
type Group struct {
	app         *App
	prefix      string
	middlewares []Middleware
}

// Use adds middlewares to this group.
func (g *Group) Use(middlewares ...Middleware) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

// Group creates a sub-group.
func (g *Group) Group(prefix string, middlewares ...Middleware) *Group {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(g.middlewares))
	copy(mws, g.middlewares)
	mws = append(mws, middlewares...)

	return &Group{
		app:         g.app,
		prefix:      g.prefix + prefix,
		middlewares: mws,
	}
}

// Route creates a RouteBuilder.
func (g *Group) Route(path string, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(g.middlewares))
	copy(mws, g.middlewares)
	mws = append(mws, middlewares...)

	return &RouteBuilder{
		app:         g.app,
		path:        g.prefix + path,
		middlewares: mws,
	}
}

// GET registers a GET handler.
func (g *Group) GET(path string, h Handler, middlewares ...Middleware) *Group {
	fullPath := g.prefix + path
	mws := append(g.middlewares, middlewares...)
	handler := chainMiddlewares(h, mws...)
	g.app.mux.Get(fullPath, g.app.wrapHandler(handler))
	return g
}

// POST registers a POST handler.
func (g *Group) POST(path string, h Handler, middlewares ...Middleware) *Group {
	fullPath := g.prefix + path
	mws := append(g.middlewares, middlewares...)
	handler := chainMiddlewares(h, mws...)
	g.app.mux.Post(fullPath, g.app.wrapHandler(handler))
	return g
}

// PUT registers a PUT handler.
func (g *Group) PUT(path string, h Handler, middlewares ...Middleware) *Group {
	fullPath := g.prefix + path
	mws := append(g.middlewares, middlewares...)
	handler := chainMiddlewares(h, mws...)
	g.app.mux.Put(fullPath, g.app.wrapHandler(handler))
	return g
}

// PATCH registers a PATCH handler.
func (g *Group) PATCH(path string, h Handler, middlewares ...Middleware) *Group {
	fullPath := g.prefix + path
	mws := append(g.middlewares, middlewares...)
	handler := chainMiddlewares(h, mws...)
	g.app.mux.Patch(fullPath, g.app.wrapHandler(handler))
	return g
}

// DELETE registers a DELETE handler.
func (g *Group) DELETE(path string, h Handler, middlewares ...Middleware) *Group {
	fullPath := g.prefix + path
	mws := append(g.middlewares, middlewares...)
	handler := chainMiddlewares(h, mws...)
	g.app.mux.Delete(fullPath, g.app.wrapHandler(handler))
	return g
}

// RouteBuilder for method chaining.
type RouteBuilder struct {
	app         *App
	path        string
	middlewares []Middleware
}

// With adds middlewares to this route.
func (rb *RouteBuilder) With(middlewares ...Middleware) *RouteBuilder {
	rb.middlewares = append(rb.middlewares, middlewares...)
	return rb
}

// GET registers a GET handler.
func (rb *RouteBuilder) GET(h Handler, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(rb.middlewares))
	copy(mws, rb.middlewares)
	mws = append(mws, middlewares...)
	handler := chainMiddlewares(h, mws...)
	rb.app.mux.Get(rb.path, rb.app.wrapHandler(handler))
	return rb
}

// POST registers a POST handler.
func (rb *RouteBuilder) POST(h Handler, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(rb.middlewares))
	copy(mws, rb.middlewares)
	mws = append(mws, middlewares...)
	handler := chainMiddlewares(h, mws...)
	rb.app.mux.Post(rb.path, rb.app.wrapHandler(handler))
	return rb
}

// PUT registers a PUT handler.
func (rb *RouteBuilder) PUT(h Handler, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(rb.middlewares))
	copy(mws, rb.middlewares)
	mws = append(mws, middlewares...)
	handler := chainMiddlewares(h, mws...)
	rb.app.mux.Put(rb.path, rb.app.wrapHandler(handler))
	return rb
}

// PATCH registers a PATCH handler.
func (rb *RouteBuilder) PATCH(h Handler, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(rb.middlewares))
	copy(mws, rb.middlewares)
	mws = append(mws, middlewares...)
	handler := chainMiddlewares(h, mws...)
	rb.app.mux.Patch(rb.path, rb.app.wrapHandler(handler))
	return rb
}

// DELETE registers a DELETE handler.
func (rb *RouteBuilder) DELETE(h Handler, middlewares ...Middleware) *RouteBuilder {
	// Copy slice to avoid sharing underlying array
	mws := make([]Middleware, len(rb.middlewares))
	copy(mws, rb.middlewares)
	mws = append(mws, middlewares...)
	handler := chainMiddlewares(h, mws...)
	rb.app.mux.Delete(rb.path, rb.app.wrapHandler(handler))
	return rb
}

// Group creates a sub-route.
func (rb *RouteBuilder) Group(subPath string, middlewares ...Middleware) *RouteBuilder {
	return &RouteBuilder{
		app:         rb.app,
		path:        rb.path + subPath,
		middlewares: append(rb.middlewares, middlewares...),
	}
}
