package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	promlog "github.com/prometheus/common/log"
)

type Server struct {
	logger *log.Logger
	mux    *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		logger: promlog.NewErrorLogger(),
		mux:    http.NewServeMux(),
	}

	s.mux.HandleFunc("/healthz", s.healthz)
	s.mux.Handle("/metrics", promhttp.Handler())

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
