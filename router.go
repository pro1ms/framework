package framework

import (
	"net/http"
)

// Router handles HTTP routing
type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// registerRoute registers a route with a specific HTTP method
func (r *Router) registerRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}

// Get registers a GET route
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.registerRoute(http.MethodGet, path, handler)
}

// Post registers a POST route
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.registerRoute(http.MethodPost, path, handler)
}

// Put registers a PUT route
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.registerRoute(http.MethodPut, path, handler)
}

// Delete registers a DELETE route
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.registerRoute(http.MethodDelete, path, handler)
}

// Patch registers a PATCH route
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.registerRoute(http.MethodPatch, path, handler)
}

// HandleFunc registers a handler function for the given pattern
// This method is required to implement the ServeMux interface from generated code
// Pattern format: "METHOD /path" (e.g., "GET /hello") or just "/path"
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// Parse pattern to extract method and path
	// Pattern format from generated code: "METHOD /path"
	method := http.MethodGet
	path := pattern

	// Check if pattern starts with HTTP method (e.g., "GET ", "POST ", etc.)
	if len(pattern) > 4 {
		spaceIdx := -1
		for i := 0; i < len(pattern) && i < 10; i++ {
			if pattern[i] == ' ' {
				spaceIdx = i
				break
			}
		}

		if spaceIdx > 0 {
			method = pattern[:spaceIdx]
			path = pattern[spaceIdx+1:]
		}
	}

	r.registerRoute(method, path, handler)
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// Check if path exists
	if methods, exists := r.routes[path]; exists {
		// Check if method exists for this path
		if handler, exists := methods[method]; exists {
			handler(w, req)
			return
		}
		// Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Path not found
	http.NotFound(w, req)
}
