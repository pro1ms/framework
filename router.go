package framework

import (
	"encoding/json"
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// Router handles HTTP routing
type Router struct {
	routes      map[string]map[string]http.HandlerFunc
	middlewares []Middleware
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) AddMiddleware(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var finalHandler http.Handler = http.HandlerFunc(r.dispatch)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}
	finalHandler.ServeHTTP(w, req)
}

// ServeHTTP implements http.Handler interface
func (r *Router) dispatch(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PANIC RECOVERED: %v\n", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			res := struct {
				Code    string
				Message string
			}{
				Code:    "internal_server_error",
				Message: "A critical error occurred on the server",
			}
			_ = json.NewEncoder(w).Encode(res)
		}
	}()

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
