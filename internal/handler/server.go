// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"sandbox/internal/config"
	"sandbox/internal/model"
	"sandbox/internal/service"
	"sandbox/internal/store"
	"sandbox/pkg/httpx"
	"sandbox/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerSandboxRoutes(mux)
	s.registerRuntimeRoutes(mux)
	s.registerExecutionTaskRoutes(mux)
	s.registerSecurityPolicyRoutes(mux)
	s.registerResourceLimitRoutes(mux)
	s.registerTemplateRoutes(mux)
	s.registerSubmissionRoutes(mux)
	s.registerScheduleRoutes(mux)
	s.registerExecutionLogRoutes(mux)
	s.registerAuditLogRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerBatchRoutes(mux)
	s.registerExportRoutes(mux)
	s.registerEnvironmentVariableRoutes(mux)
	s.registerNetworkPolicyRoutes(mux)
	s.registerVolumeMountRoutes(mux)
	s.registerArtifactRoutes(mux)
	s.registerWebhookRoutes(mux)
	s.registerNotificationRoutes(mux)
	s.registerHealthRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.apiKeyMiddleware(s.rateLimitMiddleware(s.loggingMiddleware(s.recoveryMiddleware(corsMiddleware(requestIDMiddleware(mux))))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/style.css" || r.URL.Path == "/app.js" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.APIKey {
			httpx.Unauthorized(w, "API Key 无效")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ipLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
}

func newIPLimiter(limit int) *ipLimiter {
	return &ipLimiter{requests: make(map[string][]time.Time), limit: limit}
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	window := now.Add(-1 * time.Minute)
	list := l.requests[ip]
	filtered := make([]time.Time, 0, len(list))
	for _, t := range list {
		if t.After(window) {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= l.limit {
		l.requests[ip] = filtered
		return false
	}
	filtered = append(filtered, now)
	l.requests[ip] = filtered
	return true
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := newIPLimiter(s.cfg.RateLimit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		if host == "" {
			host = r.RemoteAddr
		}
		if !limiter.allow(host) {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
