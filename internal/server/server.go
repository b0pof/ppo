package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/b0pof/ppo/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(r http.Handler, cfg config.ServerConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			ReadHeaderTimeout: 3 * time.Second,
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           r,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
