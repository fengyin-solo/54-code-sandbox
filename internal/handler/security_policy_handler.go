package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerSecurityPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/policies", s.createSecurityPolicy)
	mux.HandleFunc("GET /api/policies", s.listSecurityPolicies)
	mux.HandleFunc("GET /api/policies/{id}", s.getSecurityPolicy)
	mux.HandleFunc("PUT /api/policies/{id}", s.updateSecurityPolicy)
	mux.HandleFunc("DELETE /api/policies/{id}", s.deleteSecurityPolicy)
}

type createSecurityPolicyRequest struct {
	Name             string   `json:"name"`
	AllowNetwork     bool     `json:"allow_network"`
	AllowFileSystem  bool     `json:"allow_file_system"`
	MaxMemoryMB      int      `json:"max_memory_mb"`
	MaxCPUMs         int      `json:"max_cpu_ms"`
	MaxTimeoutMs     int      `json:"max_timeout_ms"`
	BlockedSyscalls  []string `json:"blocked_syscalls"`
	Status           string   `json:"status"`
}

func (s *Server) createSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	var req createSecurityPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateSecurityPolicy(model.SecurityPolicy{
		Name:            req.Name,
		AllowNetwork:    req.AllowNetwork,
		AllowFileSystem: req.AllowFileSystem,
		MaxMemoryMB:     req.MaxMemoryMB,
		MaxCPUMs:        req.MaxCPUMs,
		MaxTimeoutMs:    req.MaxTimeoutMs,
		BlockedSyscalls: req.BlockedSyscalls,
		Status:          req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listSecurityPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SecurityPolicyFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSecurityPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetSecurityPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updateSecurityPolicyRequest struct {
	Name             string   `json:"name"`
	AllowNetwork     bool     `json:"allow_network"`
	AllowFileSystem  bool     `json:"allow_file_system"`
	MaxMemoryMB      int      `json:"max_memory_mb"`
	MaxCPUMs         int      `json:"max_cpu_ms"`
	MaxTimeoutMs     int      `json:"max_timeout_ms"`
	BlockedSyscalls  []string `json:"blocked_syscalls"`
	Status           string   `json:"status"`
}

func (s *Server) updateSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSecurityPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateSecurityPolicy(id, model.SecurityPolicy{
		Name:            req.Name,
		AllowNetwork:    req.AllowNetwork,
		AllowFileSystem: req.AllowFileSystem,
		MaxMemoryMB:     req.MaxMemoryMB,
		MaxCPUMs:        req.MaxCPUMs,
		MaxTimeoutMs:    req.MaxTimeoutMs,
		BlockedSyscalls: req.BlockedSyscalls,
		Status:          req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteSecurityPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSecurityPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
