package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerEnvironmentVariableRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/envs", s.createEnvironmentVariable)
	mux.HandleFunc("GET /api/envs", s.listEnvironmentVariables)
	mux.HandleFunc("GET /api/envs/{id}", s.getEnvironmentVariable)
	mux.HandleFunc("PUT /api/envs/{id}", s.updateEnvironmentVariable)
	mux.HandleFunc("DELETE /api/envs/{id}", s.deleteEnvironmentVariable)
}

type createEnvironmentVariableRequest struct {
	SandboxID string `json:"sandbox_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Status    string `json:"status"`
}

func (s *Server) createEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	var req createEnvironmentVariableRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEnvironmentVariable(model.EnvironmentVariable{
		SandboxID: req.SandboxID,
		Key:       req.Key,
		Value:     req.Value,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEnvironmentVariables(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EnvironmentVariableFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Key:       r.URL.Query().Get("key"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListEnvironmentVariables(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEnvironmentVariable(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateEnvironmentVariableRequest struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Status string `json:"status"`
}

func (s *Server) updateEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEnvironmentVariableRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEnvironmentVariable(id, model.EnvironmentVariable{
		Key:    req.Key,
		Value:  req.Value,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEnvironmentVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEnvironmentVariable(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
