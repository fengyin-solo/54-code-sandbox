package service

import (
	"runtime"
	"time"
)

type HealthInfo struct {
	Status     string `json:"status"`
	Goroutines int    `json:"goroutines"`
	MemoryMB   float64 `json:"memory_mb"`
}

func (s *Service) Health() *HealthInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &HealthInfo{
		Status:     "healthy",
		Goroutines: runtime.NumGoroutine(),
		MemoryMB:   float64(m.Alloc) / 1024 / 1024,
	}
}

func (s *Service) Uptime(start time.Time) time.Duration {
	return time.Since(start)
}
