package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerResourceLimitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/limits", s.createResourceLimit)
	mux.HandleFunc("GET /api/limits", s.listResourceLimits)
	mux.HandleFunc("GET /api/limits/{id}", s.getResourceLimit)
	mux.HandleFunc("PUT /api/limits/{id}", s.updateResourceLimit)
	mux.HandleFunc("DELETE /api/limits/{id}", s.deleteResourceLimit)
}

type createResourceLimitRequest struct {
	SandboxID    string  `json:"sandbox_id"`
	CPUQuota     float64 `json:"cpu_quota"`
	MemoryMB     int     `json:"memory_mb"`
	DiskMB       int     `json:"disk_mb"`
	MaxProcesses int     `json:"max_processes"`
	Status       string  `json:"status"`
}

func (s *Server) createResourceLimit(w http.ResponseWriter, r *http.Request) {
	var req createResourceLimitRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.CreateResourceLimit(model.ResourceLimit{
		SandboxID:    req.SandboxID,
		CPUQuota:     req.CPUQuota,
		MemoryMB:     req.MemoryMB,
		DiskMB:       req.DiskMB,
		MaxProcesses: req.MaxProcesses,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, res)
}

func (s *Server) listResourceLimits(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ResourceLimitFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListResourceLimits(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getResourceLimit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := s.svc.GetResourceLimit(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

type updateResourceLimitRequest struct {
	CPUQuota     float64 `json:"cpu_quota"`
	MemoryMB     int     `json:"memory_mb"`
	DiskMB       int     `json:"disk_mb"`
	MaxProcesses int     `json:"max_processes"`
	Status       string  `json:"status"`
}

func (s *Server) updateResourceLimit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateResourceLimitRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.UpdateResourceLimit(id, model.ResourceLimit{
		CPUQuota:     req.CPUQuota,
		MemoryMB:     req.MemoryMB,
		DiskMB:       req.DiskMB,
		MaxProcesses: req.MaxProcesses,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) deleteResourceLimit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteResourceLimit(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
