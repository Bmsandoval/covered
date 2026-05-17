package api

import (
	"fmt"
	"net/http"

	"github.com/Bmsandoval/covered/internal/app"
	"github.com/Bmsandoval/covered/internal/config"
)

// Server is the HTTP transport layer.
type Server struct {
	cfg  config.Config
	deps app.Deps
	mux  *http.ServeMux
}

// NewServer builds routes. Handlers use deps when features land.
func NewServer(cfg config.Config, deps app.Deps) *Server {
	s := &Server{cfg: cfg, deps: deps, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /health", handleHealth)
	return s
}

// Handler returns the root HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// Addr is the listen address for the configured port.
func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.cfg.AppPort)
}
