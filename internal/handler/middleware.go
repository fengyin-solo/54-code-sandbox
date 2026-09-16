package handler

import (
	"crypto/subtle"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"sandbox/pkg/httpx"
	"sandbox/pkg/logger"
)

// loggingMiddleware 记录请求方法与耗时。
func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// recoveryMiddleware 捕获 panic 并返回 500。
func recoveryMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// apiKeyMiddleware 校验 API Key。
func apiKeyMiddleware(next http.Handler, cfgKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfgKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		if isStaticAsset(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if subtle.ConstantTimeCompare([]byte(key), []byte(cfgKey)) != 1 {
			httpx.Unauthorized(w, "API Key 无效")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isStaticAsset(path string) bool {
	switch path {
	case "/", "/index.html", "/style.css", "/app.js":
		return true
	}
	return false
}

// ipRateLimiter 基于 IP 的滑动窗口限流器。
type ipRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)
	list := l.requests[ip]
	filtered := make([]time.Time, 0, len(list))
	for _, t := range list {
		if t.After(cutoff) {
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

func (l *ipRateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := time.Now().Add(-l.window)
	for ip, list := range l.requests {
		filtered := make([]time.Time, 0, len(list))
		for _, t := range list {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(l.requests, ip)
		} else {
			l.requests[ip] = filtered
		}
	}
}

// rateLimitMiddleware IP 限流中间件。
func rateLimitMiddleware(next http.Handler, limit int) http.Handler {
	limiter := newIPRateLimiter(limit, time.Minute)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanup()
		}
	}()
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

// corsMiddleware 添加跨域响应头。
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware 为每个请求附加唯一请求 ID（如有客户端传入则复用）。
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = generateRequestID()
		}
		w.Header().Set("X-Request-ID", rid)
		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
