package server

import (
	"context"
	"net/http"
)

const (
	port = ":8080"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(r http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    port,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
