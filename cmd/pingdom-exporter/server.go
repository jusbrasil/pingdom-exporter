package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the object that implements the HTTP server for the exporter.
type Server struct {
	mux *http.ServeMux
}

// NewServer returns a new HTTP server for exposing Prometheus metrics.
func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("/healthz", s.healthz)
	s.mux.Handle("/metrics", promhttp.Handler())

	return s
}

// ServeHTTP handles incoming HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
