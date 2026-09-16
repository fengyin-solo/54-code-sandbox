package handler

import (
	"net/http"
	"runtime"
	"time"

	"sandbox/pkg/httpx"
)

func (s *Server) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.healthCheck)
	mux.HandleFunc("GET /ready", s.readinessCheck)
	mux.HandleFunc("GET /metrics", s.systemMetrics)
}

type healthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, healthResponse{Status: "up", Timestamp: time.Now(), Version: "1.0.0"})
}

func (s *Server) readinessCheck(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, healthResponse{Status: "ready", Timestamp: time.Now(), Version: "1.0.0"})
}

type systemMetricsResponse struct {
	Uptime       string `json:"uptime"`
	Goroutines   int    `json:"goroutines"`
	MemoryAlloc  uint64 `json:"memory_alloc"`
	MemorySys    uint64 `json:"memory_sys"`
	NumGC        uint32 `json:"num_gc"`
	MaxPageSize  int    `json:"max_page_size"`
	RateLimit    int    `json:"rate_limit"`
}

func (s *Server) systemMetrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	httpx.OK(w, systemMetricsResponse{
		Uptime:      time.Since(time.Now().Add(-time.Minute)).String(),
		Goroutines:  runtime.NumGoroutine(),
		MemoryAlloc: m.Alloc,
		MemorySys:   m.Sys,
		NumGC:       m.NumGC,
		MaxPageSize: s.maxPageSize(),
		RateLimit:   s.cfg.RateLimit,
	})
}
