package handler

import (
	"net/http"

	"sandbox/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
	mux.HandleFunc("GET /api/stats/languages", s.getLanguageDistribution)
	mux.HandleFunc("GET /api/stats/sandbox-volume", s.getSandboxExecutionVolume)
	mux.HandleFunc("GET /api/stats/timeout-rate", s.getTimeoutRate)
	mux.HandleFunc("GET /api/stats/average-duration", s.getAverageDuration)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetOverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) getLanguageDistribution(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.GetLanguageDistribution()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, data)
}

func (s *Server) getSandboxExecutionVolume(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.GetSandboxExecutionVolume()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, data)
}

func (s *Server) getTimeoutRate(w http.ResponseWriter, r *http.Request) {
	rate, err := s.svc.GetTimeoutRate()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]float64{"timeout_rate": rate})
}

func (s *Server) getAverageDuration(w http.ResponseWriter, r *http.Request) {
	avg, err := s.svc.GetAverageDuration()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]float64{"average_duration_ms": avg})
}
