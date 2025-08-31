package server

import (
	"net/http"
)

// Middleware type as before
type Middleware func(http.Handler) http.Handler

// App struct to hold our routes and middleware
type Server struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewApp creates and returns a new App with an initialized ServeMux and middleware slice
func NewServer() *Server {
	return &Server{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
	}
}

// Use adds middleware to the chain
func (s *Server) Use(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

// Handle registers a handler for a specific route, applying all middleware
func (s *Server) Handle(pattern string, handler http.Handler) {
	finalHandler := handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		finalHandler = s.middlewares[i](finalHandler)
	}
	s.mux.Handle(pattern, finalHandler)
}

// ListenAndServe starts the application server
func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.mux)
}
