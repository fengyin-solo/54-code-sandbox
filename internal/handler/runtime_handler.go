package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerRuntimeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/runtimes", s.createRuntime)
	mux.HandleFunc("GET /api/runtimes", s.listRuntimes)
	mux.HandleFunc("GET /api/runtimes/{id}", s.getRuntime)
	mux.HandleFunc("PUT /api/runtimes/{id}", s.updateRuntime)
	mux.HandleFunc("DELETE /api/runtimes/{id}", s.deleteRuntime)
}

type createRuntimeRequest struct {
	SandboxID     string  `json:"sandbox_id"`
	Name          string  `json:"name"`
	Version       string  `json:"version"`
	Command       string  `json:"command"`
	MemoryLimitMB int     `json:"memory_limit_mb"`
	CPULimit      float64 `json:"cpu_limit"`
	TimeoutMs     int     `json:"timeout_ms"`
	Status        string  `json:"status"`
}

func (s *Server) createRuntime(w http.ResponseWriter, r *http.Request) {
	var req createRuntimeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rt, err := s.svc.CreateRuntime(model.Runtime{
		SandboxID:     req.SandboxID,
		Name:          req.Name,
		Version:       req.Version,
		Command:       req.Command,
		MemoryLimitMB: req.MemoryLimitMB,
		CPULimit:      req.CPULimit,
		TimeoutMs:     req.TimeoutMs,
		Status:        req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rt)
}

func (s *Server) listRuntimes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuntimeFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRuntimes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRuntime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rt, err := s.svc.GetRuntime(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rt)
}

type updateRuntimeRequest struct {
	Name          string  `json:"name"`
	Version       string  `json:"version"`
	Command       string  `json:"command"`
	MemoryLimitMB int     `json:"memory_limit_mb"`
	CPULimit      float64 `json:"cpu_limit"`
	TimeoutMs     int     `json:"timeout_ms"`
	Status        string  `json:"status"`
}

func (s *Server) updateRuntime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRuntimeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rt, err := s.svc.UpdateRuntime(id, model.Runtime{
		Name:          req.Name,
		Version:       req.Version,
		Command:       req.Command,
		MemoryLimitMB: req.MemoryLimitMB,
		CPULimit:      req.CPULimit,
		TimeoutMs:     req.TimeoutMs,
		Status:        req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rt)
}

func (s *Server) deleteRuntime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRuntime(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
