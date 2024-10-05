package pkg

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Server struct to hold registered routes
type Server struct {
	*http.ServeMux
	registeredRoutes []string
}

// HandleFunc is a custom handler function to track and register routes
func (m *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// Log the route when it's registered
	m.registeredRoutes = append(m.registeredRoutes, pattern)
	// Register the actual handler
	m.ServeMux.HandleFunc(pattern, handler)
}

// PrintRoutes prints all registered routes
func (m *Server) PrintRoutes() {
	slog.Info("Registered Routes:")
	var routes string
	for _, route := range m.registeredRoutes {
		routes += fmt.Sprintf(" - Route %s\n", route)
	}
	slog.Info(routes)
}
